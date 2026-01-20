import styles from './registration-form.module.css';
import {useEffect, useState} from 'react';
import {EmailValidator, NameValidator, NicknameValidator, PasswordValidator} from "@/app/login/validators";
import {Login, Register} from "@/api/auth";
import {Message} from "@/api/api";
import {useRouter} from "next/navigation";
import {useUserStore} from "@/context/user-store";
import MessageComponent from "@/app/components/message";
import Spinner from "@/app/components/spinner";

interface FormData {
    name: string;
    surname: string;
    patronymic: string;
    nickname: string
    email: string;
    password: string;
    confirmPassword: string;
    day: number;
    month: number;
    year: number;
    gender: 'male' | 'female';
}

const MONTHS = [
    'января',
    'февраля',
    'марта',
    'апреля',
    'мая',
    'июня',
    'июля',
    'августа',
    'сентября',
    'октября',
    'ноября',
    'декабря',
];

export default function RegistrationForm({onLoginClick}: { onLoginClick?: () => void }) {
    const currentYear = new Date().getFullYear();
    const [formData, setFormData] = useState<FormData>({
        name: '',
        surname: '',
        patronymic: '',
        nickname: '',
        email: '',
        password: '',
        confirmPassword: '',
        day: 1,
        month: 1,
        year: currentYear - 18,
        gender: 'male',
    });

    const [response, setResponse] = useState<Message | null>(null)
    const [requestSent, setRequestSent] = useState<boolean>(false);
    const router = useRouter();
    const user = useUserStore((state) => state);

    useEffect(() => {
        response && setRequestSent(false);
        if (user.user)
            router.push('/');
    }, [response, router, user]);

    const birthDate = new Date(
        formData.year,
        formData.month - 1,
        formData.day
    );
    const today = new Date();
    const ageDiff = today.getFullYear() - birthDate.getFullYear();
    const ageCheck =
        ageDiff > 18 ||
        (ageDiff === 18 &&
            (today.getMonth() > birthDate.getMonth() ||
                (today.getMonth() === birthDate.getMonth() &&
                    today.getDate() >= birthDate.getDate())));

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        Register({
            name: formData.name,
            surname: formData.surname,
            patronymic: formData.patronymic,
            nickname: formData.nickname,
            email: formData.email,
            password: formData.password,
            birthDate: birthDate
        }, setResponse);
    };

    return (
        <form className={styles.registration_form} onSubmit={handleSubmit}>
            <div className={styles.auth_options}>
                <h2>Регистрация</h2>
                <h2 onClick={onLoginClick}>Войти</h2>
            </div>
            <hr style={{width: 'calc(100% + 40px)', margin: '10px -20px'}}/>
            <input
                className={styles.input_field}
                placeholder="Имя"
                onInput={NameValidator} required minLength={2}
                value={formData.name}
                onChange={(e) =>
                    setFormData({...formData, name: e.target.value})
                }
            />
            <input
                className={styles.input_field}
                placeholder="Фамилия"
                onInput={NameValidator} required minLength={2}
                value={formData.surname}
                onChange={(e) =>
                    setFormData({...formData, surname: e.target.value})
                }
            />
            <input
                className={styles.input_field}
                placeholder="Отчество"
                onInput={NameValidator} required minLength={2}
                value={formData.patronymic}
                onChange={(e) =>
                    setFormData({...formData, patronymic: e.target.value})
                }
            />
            <input
                className={styles.input_field}
                placeholder="Имя профиля"
                onInput={NicknameValidator} required minLength={2}
                value={formData.nickname}
                onChange={(e) =>
                    setFormData({...formData, nickname: e.target.value})
                }
            />
            <input
                className={styles.input_field}
                placeholder="Email"
                onInput={EmailValidator} required
                value={formData.email}
                onChange={(e) =>
                    setFormData({...formData, email: e.target.value})
                }
            />
            <input
                className={styles.input_field}
                placeholder="Пароль"
                onInput={PasswordValidator} required minLength={8}
                type="password"
                value={formData.password}
                onChange={(e) =>
                    setFormData({...formData, password: e.target.value})
                }
            />
            <input
                className={styles.input_field}
                placeholder="Повторите пароль"
                onInput={e =>
                    (PasswordValidator(e), e.currentTarget.checkValidity()
                    && e.currentTarget.setCustomValidity(e.currentTarget.value === formData.password ? '' : 'Пароли не совпадают'),
                        e.currentTarget.reportValidity())} required
                minLength={8}
                type="password"
                value={formData.confirmPassword}
                onChange={(e) =>
                    setFormData({...formData, confirmPassword: e.target.value})
                }
            />
            <p className={styles.field_label}>Дата рождения</p>
            <div className={styles.birth_date_selects}>
                <select
                    value={formData.day}
                    onChange={(e) =>
                        setFormData({...formData, day: parseInt(e.target.value)})
                    }
                >
                    {Array.from({length: monthDays(formData.year, formData.month)}, (_, i) => (
                        <option key={i + 1} value={i + 1}>
                            {i + 1}
                        </option>
                    ))}
                </select>
                <select
                    value={formData.month}
                    onChange={(e) =>
                        setFormData({...formData, month: parseInt(e.target.value)})
                    }
                >
                    {MONTHS.map((month, index) => (
                        <option key={index + 1} value={index + 1}>
                            {month}
                        </option>
                    ))}
                </select>
                <select
                    value={formData.year}
                    onChange={(e) =>
                        setFormData({...formData, year: parseInt(e.target.value)})
                    }
                >
                    {Array.from({length: 100 - 17}, (_, i) => currentYear + i - 100)
                        .reverse()
                        .map((year) => (
                            <option key={year} value={year}>
                                {year}
                            </option>
                        ))}
                </select>
            </div>
            <p className={styles.field_label}>Пол</p>
            <div className={styles.gender_group}>
                <label>
                    <input
                        type="radio"
                        name="gender"
                        value="male"
                        checked={formData.gender === 'male'}
                        onChange={() => setFormData({...formData, gender: 'male'})}
                    />
                    Мужской
                </label>
                <label>
                    <input
                        type="radio"
                        name="gender"
                        value="female"
                        checked={formData.gender === 'female'}
                        onChange={() => setFormData({...formData, gender: 'female'})}
                    />
                    Женский
                </label>
                {/*<label>*/}
                {/*    <input*/}
                {/*        type="radio"*/}
                {/*        name="gender"*/}
                {/*        value="not-selected"*/}
                {/*        checked={formData.gender === 'not-selected'}*/}
                {/*        onChange={() =>*/}
                {/*            setFormData({...formData, gender: 'not-selected'})*/}
                {/*        }*/}
                {/*    />*/}
                {/*    Не выбран*/}
                {/*</label>*/}
            </div>
            <hr style={{width: 'calc(100% + 40px)', margin: '10px -20px'}}/>
            <MessageComponent message={response}/>
            <button disabled={requestSent} className={styles.submit_button}>{requestSent &&
                <Spinner size={30} style={{margin: "-11px 0 -11px -32px", paddingRight: "32px"}}/>}Зарегистрироваться
            </button>
            <p className={styles.agreement_text}>
                Нажимая кнопку «Зарегистрироваться» вы принимаете условия
                пользовательского соглашения, условия публичной оферты и условий
                политики конфиденциальности и даете согласие на обработку
                персональных данных
            </p>
        </form>
    );
}

function monthDays(year: number, monthIndex: number) {
    var d = new Date(year, monthIndex, 0);
    return d.getDate();
}