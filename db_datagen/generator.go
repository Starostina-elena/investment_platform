package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	NumUsers         int
	NumOrganizations int
	NumProjects      int
	NumTags          int
}

type Generator struct {
	conn             *pgx.Conn
	config           Config
	userIDs          []int
	organizationIDs  []int
	projectIDs       []int
	tagIDs           []int
	orgToPhysAccount map[int]int
	orgToJurAccount  map[int]int
	orgToIpAccount   map[int]int
}

func NewGenerator(conn *pgx.Conn, config Config) *Generator {
	gofakeit.Seed(time.Now().UnixNano())
	return &Generator{
		conn:             conn,
		config:           config,
		orgToPhysAccount: make(map[int]int),
		orgToJurAccount:  make(map[int]int),
		orgToIpAccount:   make(map[int]int),
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (g *Generator) Run() error {
	fmt.Println("Начинаем генерацию тестовых данных. За работу, товарищи!")

	if err := g.generateUsers(); err != nil {
		log.Printf(">> %v", err)
		return err
	}
	if err := g.generateOrganizations(); err != nil {
		return err
	}
	if err := g.generateOrgAccounts(); err != nil {
		return err
	}
	if err := g.generateTags(); err != nil {
		return err
	}
	if err := g.generateProjects(); err != nil {
		return err
	}
	if err := g.generateUserRights(); err != nil {
		return err
	}
	if err := g.generateComments(); err != nil {
		return err
	}
	if err := g.generateTransactions(); err != nil {
		return err
	}

	fmt.Println("Перепись населения и национализация завершены успешно! Пятилетка - в три года!")
	return nil
}

func (g *Generator) generateProjects() error {
	fmt.Printf("Запускаем %d великих строек (проектов)...\n", g.config.NumProjects)
	if len(g.organizationIDs) == 0 || len(g.tagIDs) == 0 {
		return fmt.Errorf("нет организаций или тэгов для создания проектов!")
	}

	tx, err := g.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	projectSQL := `INSERT INTO projects (name, creator_id, quick_peek, content, wanted_money, duration_days) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	projectTagSQL := `INSERT INTO project_tags (project_id, tag_id) VALUES ($1, $2)`

	for i := 0; i < g.config.NumProjects; i++ {
		name := fmt.Sprintf("%s %s", gofakeit.RandomString(projectAdjectives), gofakeit.RandomString(projectNouns))
		creatorID := g.organizationIDs[rand.Intn(len(g.organizationIDs))]
		quickPeek := gofakeit.Sentence(20)
		content := gofakeit.Paragraph(5, 10, 50, "\n")
		wantedMoney := gofakeit.Price(50000, 5000000)
		duration := gofakeit.Number(30, 365)

		var projectID int
		err := tx.QueryRow(context.Background(), projectSQL,
			name, creatorID, quickPeek, content, wantedMoney, duration,
		).Scan(&projectID)
		if err != nil {
			return err
		}
		g.projectIDs = append(g.projectIDs, projectID)

		numTagsForProject := rand.Intn(3) + 1
		rand.Shuffle(len(g.tagIDs), func(i, j int) { g.tagIDs[i], g.tagIDs[j] = g.tagIDs[j], g.tagIDs[i] })

		for j := 0; j < numTagsForProject; j++ {
			tagID := g.tagIDs[j]
			if _, err := tx.Exec(context.Background(), projectTagSQL, projectID, tagID); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}

	fmt.Println("Проекты успешно созданы и категоризированы.")
	return nil
}

func (g *Generator) generateTags() error {
	fmt.Printf("Создаем %d тэгов для категоризации проектов...\n", g.config.NumTags)

	batch := &pgx.Batch{}
	sql := `INSERT INTO tags (name, description) VALUES ($1, $2)`

	uniqueTags := make(map[string]bool)
	for len(uniqueTags) < g.config.NumTags {
		uniqueTags[gofakeit.RandomString(tags)] = true
	}

	for tagName := range uniqueTags {
		batch.Queue(sql, tagName, gofakeit.Sentence(10))
	}

	br := g.conn.SendBatch(context.Background(), batch)
	for i := 0; i < len(uniqueTags); i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	if err := br.Close(); err != nil {
		return err
	}

	rows, err := g.conn.Query(context.Background(), "SELECT id FROM tags")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		g.tagIDs = append(g.tagIDs, id)
	}

	fmt.Println("Тэги созданы.")
	return nil
}

type UserData struct {
	Name         string
	Surname      string
	Patronymic   string
	Nickname     string
	Email        string
	PasswordHash string
	Balance      float64
	IsAdmin      bool
}

type UserCopySource struct {
	users []UserData
	idx   int
}

func (ucs *UserCopySource) Next() bool {
	ucs.idx++
	return ucs.idx-1 < len(ucs.users)
}

func (ucs *UserCopySource) Values() ([]interface{}, error) {
	user := ucs.users[ucs.idx-1]
	return []interface{}{
		user.Name, user.Surname, user.Patronymic, user.Nickname,
		user.Email, user.PasswordHash, user.Balance, user.IsAdmin,
	}, nil
}

func (ucs *UserCopySource) Err() error {
	return nil
}

func (g *Generator) generateUsers() error {
	fmt.Printf("Генерируем %d пользователей...\n", g.config.NumUsers)

	usersData := make([]UserData, g.config.NumUsers)
	for i := 0; i < g.config.NumUsers; i++ {
		isMale := gofakeit.Bool()
		var name, surname, patronymic string

		if isMale {
			name = gofakeit.RandomString(maleNames)
			surname = gofakeit.RandomString(surnamesMale)
			patronymic = gofakeit.RandomString(patronymicsMale)
		} else {
			name = gofakeit.RandomString(femaleNames)
			surname = gofakeit.RandomString(surnamesFemale)
			patronymic = gofakeit.RandomString(patronymicsFemale)
		}

		passwordHash, _ := hashPassword("password123")

		usersData[i] = UserData{
			Name:         name,
			Surname:      surname,
			Patronymic:   patronymic,
			Nickname:     fmt.Sprintf("%s_%d", gofakeit.Username(), i),
			Email:        gofakeit.Email(),
			PasswordHash: passwordHash,
			Balance:      gofakeit.Price(100, 10000),
			IsAdmin:      gofakeit.Number(0, 100) < 1,
		}
	}

	source := &UserCopySource{users: usersData, idx: 0}

	tx, err := g.conn.Begin(context.Background())
	if err != nil {
		log.Printf(">> %v", err)
		return err
	}
	defer tx.Rollback(context.Background())

	tableName := pgx.Identifier{"users"}
	columnNames := []string{
		"name", "surname", "patronymic", "nickname", "email",
		"password_hash", "balance", "is_admin",
	}

	copyCount, err := tx.CopyFrom(context.Background(), tableName, columnNames, source)
	if err != nil {
		return err
	}

	if int(copyCount) != g.config.NumUsers {
		return fmt.Errorf("ожидалось вставить %d юзеров, а вставилось %d", g.config.NumUsers, copyCount)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}

	rows, err := g.conn.Query(context.Background(), "SELECT id FROM USERS ORDER BY id DESC LIMIT $1", g.config.NumUsers)
	if err != nil {
		return err
	}
	defer rows.Close()

	g.userIDs = g.userIDs[:0]
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		g.userIDs = append(g.userIDs, id)
	}
	for i, j := 0, len(g.userIDs)-1; i < j; i, j = i+1, j-1 {
		g.userIDs[i], g.userIDs[j] = g.userIDs[j], g.userIDs[i]
	}

	fmt.Println("Пользователи успешно созданы.")
	return nil
}

func (g *Generator) generateOrgAccounts() error {
	fmt.Println("Создаем бюрократический аппарат: счета для организаций...")

	rows, err := g.conn.Query(context.Background(), "SELECT id, type FROM ORGANIZATIONS")
	if err != nil {
		return err
	}
	defer rows.Close()

	orgTypes := make(map[int]string)
	for rows.Next() {
		var id int
		var orgType string
		if err := rows.Scan(&id, &orgType); err != nil {
			return err
		}
		orgTypes[id] = orgType
	}

	for orgID, orgType := range orgTypes {
		var accID int
		var err error

		switch orgType {
		case "phys":
			accID, err = g.createPhysicalFaceAccount()
			if err != nil {
				return err
			}
			g.orgToPhysAccount[orgID] = accID
		case "jur":
			accID, err = g.createJuridicalFaceAccount()
			if err != nil {
				return err
			}
			g.orgToJurAccount[orgID] = accID
		case "ip":
			accID, err = g.createIpAccount()
			if err != nil {
				return err
			}
			g.orgToIpAccount[orgID] = accID
		}

		_, err = g.conn.Exec(context.Background(),
			"UPDATE ORGANIZATIONS SET org_type_id = $1 WHERE id = $2", accID, orgID)
		if err != nil {
			return err
		}
	}

	fmt.Println("Счета успешно созданы и привязаны.")
	return nil
}

func (g *Generator) createPhysicalFaceAccount() (int, error) {
	var id int
	sql := `INSERT INTO physical_face_project_account (BIC, checking_account, correspondent_account, FIO, INN, pasport_series, pasport_number, pasport_givenby, registration_address, post_address, pasport_page_with_photo_path, pasport_page_with_propiska_path, svid_o_postanovke_na_uchet_phys_litsa_path) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`

	fio := gofakeit.Name()
	address := gofakeit.Address()

	err := g.conn.QueryRow(context.Background(), sql,
		gofakeit.Number(100000000, 999999999), gofakeit.Number(1000000000, 9999999999), gofakeit.Number(1000000000, 9999999999),
		fio, gofakeit.Number(100000000000, 999999999999), gofakeit.Number(1000, 9999), gofakeit.Number(100000, 999999),
		fmt.Sprintf("ОВД %s", gofakeit.City()), address.Address, address.Address,
		gofakeit.ImageURL(200, 300), gofakeit.ImageURL(200, 300), gofakeit.ImageURL(200, 300),
	).Scan(&id)
	return id, err
}

func (g *Generator) createJuridicalFaceAccount() (int, error) {
	var id int
	sql := `INSERT INTO juridical_face_project_accout (acts_on_base, position, BIC, checking_account, correspondent_account, full_organisation_name, short_organisation_name, INN, OGRN, KPP, jur_address, fact_address, post_address, svid_o_registratsii_jur_litsa_path, svid_o_postanovke_na_nalog_uchet_path, protocol_o_nasznachenii_litsa_path, USN_path, ustav_path) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING id`

	shortName := fmt.Sprintf("ООО '%s'", gofakeit.Company())
	address := gofakeit.Address()

	err := g.conn.QueryRow(context.Background(), sql,
		"Устав", gofakeit.JobTitle(), gofakeit.Number(100000000, 999999999), gofakeit.Number(10000000000000, 99999999999999), gofakeit.Number(10000000000000, 99999999999999),
		fmt.Sprintf("Общество с ограниченной ответственностью \"%s\"", gofakeit.Company()), shortName, gofakeit.Number(1000000000, 9999999999), gofakeit.Number(1000000000000, 9999999999999),
		fmt.Sprintf("%d", gofakeit.Number(100000000, 999999999)), address.Address, address.Address, address.Address,
		gofakeit.ImageURL(300, 400), gofakeit.ImageURL(300, 400), gofakeit.ImageURL(300, 400), gofakeit.ImageURL(300, 400), gofakeit.ImageURL(300, 400),
	).Scan(&id)
	return id, err
}

func (g *Generator) createIpAccount() (int, error) {
	var id int
	sql := `INSERT INTO ip_project_account (BIC, ras_schot, kor_schot, FIO, ip_svid_serial, ip_svid_number, ip_svid_givenby, INN, OGRN, jur_address, fact_address, post_address, svid_o_postanovke_na_nalog_uchet_path, ip_pasport_photo_page_path, ip_pasport_propiska_path, USN_path, OGRNIP_path)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id`

	address := gofakeit.Address()
	fio := gofakeit.Name()

	err := g.conn.QueryRow(context.Background(), sql,
		gofakeit.Number(100000000, 999999999), gofakeit.Number(10000000000000, 99999999999999), gofakeit.Number(10000000000000, 99999999999999),
		fio, gofakeit.Number(1000, 9999), gofakeit.Number(100000, 999999), "МИФНС России", gofakeit.Number(100000000000, 999999999999), gofakeit.Number(1000000000000, 9999999999999),
		address.Address, address.Address, address.Address,
		gofakeit.ImageURL(300, 400), gofakeit.ImageURL(300, 400), gofakeit.ImageURL(300, 400), gofakeit.ImageURL(300, 400), gofakeit.ImageURL(300, 400),
	).Scan(&id)
	return id, err
}

func (g *Generator) generateOrganizations() error {
	fmt.Printf("Генерируем %d организаций...\n", g.config.NumOrganizations)
	if len(g.userIDs) == 0 {
		return fmt.Errorf("нет пользователей, чтобы назначить их владельцами организаций. Саботаж!")
	}

	tx, err := g.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	batch := &pgx.Batch{}
	sql := `
		INSERT INTO ORGANIZATIONS (name, owner, email, balance, type)
		VALUES ($1, $2, $3, $4, $5)
	`
	orgTypes := []string{"jur", "phys", "ip"}

	for i := 0; i < g.config.NumOrganizations; i++ {
		name := fmt.Sprintf("%s '%s'", gofakeit.RandomString(orgNamesPrefix), gofakeit.RandomString(orgNamesSuffix))
		ownerID := g.userIDs[rand.Intn(len(g.userIDs))]
		email := gofakeit.Email()
		balance := gofakeit.Price(10000, 1000000)
		orgType := gofakeit.RandomString(orgTypes)
		batch.Queue(sql, name, ownerID, email, balance, orgType)
	}

	br := tx.SendBatch(context.Background(), batch)
	for i := 0; i < g.config.NumOrganizations; i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("ошибка при вставке организации: %w", err)
		}
	}
	if err := br.Close(); err != nil {
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}

	rows, err := g.conn.Query(context.Background(), "SELECT id FROM ORGANIZATIONS ORDER BY id DESC LIMIT $1", g.config.NumOrganizations)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		g.organizationIDs = append(g.organizationIDs, id)
	}
	for i, j := 0, len(g.organizationIDs)-1; i < j; i, j = i+1, j-1 {
		g.organizationIDs[i], g.organizationIDs[j] = g.organizationIDs[j], g.organizationIDs[i]
	}

	fmt.Println("Организации успешно созданы.")
	return nil
}

func (g *Generator) generateUserRights() error {
	fmt.Println("Выдаем партийные билеты (права доступа в организациях)...")
	if len(g.userIDs) == 0 || len(g.organizationIDs) == 0 {
		return fmt.Errorf("нет пользователей или организаций для назначения прав!")
	}

	batch := &pgx.Batch{}
	sql := `INSERT INTO user_right_at_org (org_id, user_id, org_account_management, money_management, project_management)
	VALUES ($1, $2, $3, $4, $5)`

	insertCount := 0
	for _, userID := range g.userIDs {
		numOrgs := rand.Intn(3) + 1
		rand.Shuffle(len(g.organizationIDs), func(i, j int) {
			g.organizationIDs[i], g.organizationIDs[j] = g.organizationIDs[j], g.organizationIDs[i]
		})

		for i := 0; i < numOrgs && i < len(g.organizationIDs); i++ {
			orgID := g.organizationIDs[i]
			batch.Queue(sql, orgID, userID, gofakeit.Bool(), gofakeit.Bool(), gofakeit.Bool())
			insertCount++
		}
	}

	br := g.conn.SendBatch(context.Background(), batch)
	for i := 0; i < insertCount; i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	if err := br.Close(); err != nil {
		return err
	}

	fmt.Printf("Выдано %d удостоверений.\n", insertCount)
	return nil
}

func (g *Generator) generateComments() error {
	numComments := g.config.NumProjects * 5
	fmt.Printf("Разводим срачи в комментариях (%d штук)...\n", numComments)
	if len(g.userIDs) == 0 || len(g.projectIDs) == 0 {
		return fmt.Errorf("некому или негде комментировать!")
	}

	batch := &pgx.Batch{}
	sql := `INSERT INTO comments (body, user_id, project_id, created_at) VALUES ($1, $2, $3, $4)`

	for i := 0; i < numComments; i++ {
		userID := g.userIDs[rand.Intn(len(g.userIDs))]
		projectID := g.projectIDs[rand.Intn(len(g.projectIDs))]
		body := gofakeit.Paragraph(1, 3, 15, " ")
		createdAt := gofakeit.DateRange(time.Now().AddDate(0, -1, 0), time.Now())
		batch.Queue(sql, body, userID, projectID, createdAt)
	}

	br := g.conn.SendBatch(context.Background(), batch)
	for i := 0; i < numComments; i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	if err := br.Close(); err != nil {
		return err
	}

	fmt.Println("Комментарии успешно оставлены.")
	return nil
}

func (g *Generator) generateTransactions() error {
	numTransactions := g.config.NumUsers * 10000 // по 10 транзакций на юзера
	fmt.Printf("Запускаем денежные потоки (%d транзакций)...\n", numTransactions)

	txTypes := []string{
		"user_to_project", "user_deposit", "user_withdraw",
	}

	for i := 0; i < numTransactions; i++ {
		txType := gofakeit.RandomString(txTypes)
		amount := gofakeit.Price(100, 5000)

		tx, err := g.conn.Begin(context.Background())
		if err != nil {
			return err
		}

		var fromID, toID int
		var fromBalance, toBalance float64
		var fromTable, toTable, fromColumn, toColumn string

		switch txType {
		case "user_deposit":
			fromID = 0 // Системный источник
			toID = g.userIDs[rand.Intn(len(g.userIDs))]
			fromTable, toTable = "", "USERS"
			fromColumn, toColumn = "", "balance"
		case "user_withdraw":
			fromID = g.userIDs[rand.Intn(len(g.userIDs))]
			toID = 0 // Системный сток
			fromTable, toTable = "USERS", ""
			fromColumn, toColumn = "balance", ""
		case "user_to_project":
			fromID = g.userIDs[rand.Intn(len(g.userIDs))]
			toID = g.projectIDs[rand.Intn(len(g.projectIDs))]
			fromTable, toTable = "USERS", "projects"
			fromColumn, toColumn = "balance", "current_money"
		}

		if fromTable != "" {
			err = tx.QueryRow(context.Background(), fmt.Sprintf("SELECT %s FROM %s WHERE id = $1 FOR UPDATE", fromColumn, fromTable), fromID).Scan(&fromBalance)
			if err != nil {
				tx.Rollback(context.Background())
				continue
			}
		}
		if toTable != "" {
			err = tx.QueryRow(context.Background(), fmt.Sprintf("SELECT %s FROM %s WHERE id = $1 FOR UPDATE", toColumn, toTable), toID).Scan(&toBalance)
			if err != nil {
				tx.Rollback(context.Background())
				continue
			}
		}

		if fromTable != "" && fromBalance < amount {
			tx.Rollback(context.Background())
			continue
		}

		newFromBalance, newToBalance := fromBalance, toBalance
		if fromTable != "" {
			newFromBalance = fromBalance - amount
			_, err = tx.Exec(context.Background(), fmt.Sprintf("UPDATE %s SET %s = $1 WHERE id = $2", fromTable, fromColumn), newFromBalance, fromID)
			if err != nil {
				tx.Rollback(context.Background())
				continue
			}
		}
		if toTable != "" {
			newToBalance = toBalance + amount
			_, err = tx.Exec(context.Background(), fmt.Sprintf("UPDATE %s SET %s = $1 WHERE id = $2", toTable, toColumn), newToBalance, toID)
			if err != nil {
				tx.Rollback(context.Background())
				continue
			}
		}

		_, err = tx.Exec(context.Background(),
			`INSERT INTO transactions (from_id, reciever_id, type, amount, cum_sum_of_sender, cum_sum_of_reciever)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			fromID, toID, txType, amount, newFromBalance, newToBalance)
		if err != nil {
			tx.Rollback(context.Background())
			continue
		}

		if err := tx.Commit(context.Background()); err != nil {
			fmt.Printf("Ошибка коммита транзакции: %v\n", err)
		}
	}

	fmt.Println("Финансовые операции проведены.")
	return nil
}
