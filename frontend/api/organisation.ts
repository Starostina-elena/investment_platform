import {api, DefaultErrorHandler, Message} from "@/api/api";

// Типы лиц
export type OrgType = 'phys' | 'jur' | 'ip';

// Базовые поля
export interface Organisation {
    id: number;
    name: string;
    email: string;
    owner_id: number;
    avatar_path?: string;
    balance: number;
    org_type: OrgType;
    is_banned: boolean;
    registration_completed: boolean;
    created_at?: string;

    // Вложенные структуры (приходят с бэка)
    phys_face?: PhysFace;
    jur_face?: JurFace;
    ip_face?: IpFace;
}

export interface PhysFace {
    bic: string;
    checking_account: string;
    correspondent_account: string;
    fio: string;
    inn: string;
    passport_series: number;
    passport_number: number;
    passport_givenby: string;
    registration_address: string;
    post_address: string;
}

export interface JurFace {
    acts_on_base: string;
    position: string;
    bic: string;
    checking_account: string;
    correspondent_account: string;
    full_organisation_name: string;
    short_organisation_name: string;
    inn: string;
    ogrn: string;
    kpp: string;
    jur_address: string;
    fact_address: string;
    post_address: string;
}

export interface IpFace {
    bic: string;
    ras_schot: string;
    kor_schot: string;
    fio: string;
    ip_svid_serial: number;
    ip_svid_number: number;
    ip_svid_givenby: string;
    inn: string;
    ogrn: string;
    jur_address: string;
    fact_address: string;
    post_address: string;
}

// Получить список моих организаций
export async function GetMyOrganisations(): Promise<Organisation[]> {
    try {
        const res = await api.get('/org/my');
        return Array.isArray(res.data) ? res.data : [];
    } catch (e) {
        console.warn(e);
        return [];
    }
}

// Получить полную инфу об организации (для владельца)
export async function GetFullOrganisation(id: number): Promise<Organisation | null> {
    try {
        const res = await api.get(`/org/${id}/full`);
        return res.data;
    } catch (e) {
        console.warn(e);
        return null;
    }
}

// Создание
export async function CreateOrganisation(org: Partial<Organisation>, setMessage: (msg: Message) => void): Promise<Organisation | null> {
    try {
        const res = await api.post('/org/create', org);
        setMessage({isError: false, message: "Организация создана!"});
        return res.data;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return null;
    }
}

// Обновление
export async function UpdateOrganisation(id: number, org: Partial<Organisation>, setMessage: (msg: Message) => void): Promise<Organisation | null> {
    try {
        const res = await api.post(`/org/${id}/update`, org);
        setMessage({isError: false, message: "Данные обновлены!"});
        return res.data;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return null;
    }
}

// Аватар
export async function UploadOrgAvatar(id: number, file: File, setMessage: (msg: Message) => void): Promise<string | null> {
    try {
        const formData = new FormData();
        formData.append("avatar", file);

        const res = await api.post(`/org/${id}/avatar/upload`, formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
        setMessage({isError: false, message: "Логотип обновлен"});
        return res.data.avatar_path;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return null;
    }
}

export async function BanOrganisation(orgId: number, ban: boolean, setMessage: (msg: Message) => void): Promise<boolean> {
    try {
        await api.post(`/org/${orgId}/active?ban=${ban}`);
        setMessage({isError: false, message: ban ? "Организация заблокирована" : "Организация разблокирована"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

export async function DeleteOrgAvatar(orgId: number, setMessage: (msg: Message) => void): Promise<boolean> {
    try {
        await api.delete(`/org/${orgId}/avatar`);
        setMessage({isError: false, message: "Логотип удален"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

// Сотрудники
export interface Employee {
    org_id: number;
    user_id: number;
    user_email: string;
    nickname: string;
    org_account_management: boolean;
    money_management: boolean;
    project_management: boolean;
}

export async function GetOrganisationEmployees(orgId: number): Promise<Employee[]> {
    try {
        const res = await api.get(`/org/${orgId}/employees`);
        return Array.isArray(res.data) ? res.data : [];
    } catch (e) {
        console.warn(e);
        return [];
    }
}

export async function AddEmployee(
    orgId: number,
    userId: number,
    permissions: {
        org_account_management: boolean;
        money_management: boolean;
        project_management: boolean;
    },
    setMessage: (msg: Message) => void
): Promise<boolean> {
    try {
        await api.post(`/org/${orgId}/employees/add`, {
            user_id: userId,
            ...permissions
        });
        setMessage({isError: false, message: "Сотрудник добавлен"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

export async function UpdateEmployee(
    orgId: number,
    userId: number,
    permissions: {
        org_account_management: boolean;
        money_management: boolean;
        project_management: boolean;
    },
    setMessage: (msg: Message) => void
): Promise<boolean> {
    try {
        await api.post(`/org/${orgId}/employees/update`, {
            user_id: userId,
            ...permissions
        });
        setMessage({isError: false, message: "Права сотрудника обновлены"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

export async function DeleteEmployee(
    orgId: number,
    userId: number,
    setMessage: (msg: Message) => void
): Promise<boolean> {
    try {
        await api.delete(`/org/${orgId}/employees/${userId}/delete`);
        setMessage({isError: false, message: "Сотрудник удален"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

export async function TransferOwnership(
    orgId: number,
    newOwnerId: number,
    setMessage: (msg: Message) => void
): Promise<boolean> {
    try {
        await api.post(`/org/${orgId}/ownership/transfer/${newOwnerId}`);
        setMessage({isError: false, message: "Право владения передано"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

// Документы
export type DocType =
    | 'phys_passport_photo'
    | 'phys_passport_propiska'
    | 'phys_svid_uchet'
    | 'jur_reg_svid'
    | 'jur_svid_uchet'
    | 'jur_appointment_protocol'
    | 'jur_usn'
    | 'jur_ustav'
    | 'ip_svid_uchet'
    | 'ip_passport_photo'
    | 'ip_passport_propiska'
    | 'ip_usn'
    | 'ip_ogrnip';

export async function UploadOrgDocument(
    orgId: number,
    docType: DocType,
    file: File,
    setMessage: (msg: Message) => void
): Promise<boolean> {
    try {
        const formData = new FormData();
        formData.append('file', file);

        await api.post(`/org/${orgId}/docs/${docType}`, formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });

        setMessage({isError: false, message: "Документ загружен"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

export async function GetOrgDocument(
    orgId: number,
    docType: DocType
): Promise<{ blob: Blob; filename: string } | null> {
    try {
        const res = await api.get(`/org/${orgId}/docs/${docType}`, {
            responseType: 'blob'
        });
        
        // Получаем имя файла из заголовка Content-Disposition
        const contentDisposition = res.headers['content-disposition'];
        let filename = `${docType}`;
        
        if (contentDisposition) {
            const filenameMatch = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/);
            if (filenameMatch && filenameMatch[1]) {
                filename = filenameMatch[1].replace(/['"]/g, '');
            }
        }
        
        // Если имя файла не получили, определяем расширение по Content-Type
        if (!filename.includes('.')) {
            const contentType = res.headers['content-type'];
            const extension = getExtensionFromContentType(contentType);
            filename = `${docType}${extension}`;
        }
        
        return { blob: res.data, filename };
    } catch (e) {
        console.warn(e);
        return null;
    }
}

function getExtensionFromContentType(contentType: string): string {
    const typeMap: Record<string, string> = {
        'application/pdf': '.pdf',
        'image/jpeg': '.jpg',
        'image/jpg': '.jpg',
        'image/png': '.png',
        'image/tiff': '.tiff',
        'application/msword': '.doc',
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document': '.docx',
        'application/vnd.ms-excel': '.xls',
        'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet': '.xlsx'
    };
    
    return typeMap[contentType] || '.pdf';
}

export async function DeleteOrgDocument(
    orgId: number,
    docType: DocType,
    setMessage: (msg: Message) => void
): Promise<boolean> {
    try {
        await api.delete(`/org/${orgId}/docs/${docType}`);
        setMessage({isError: false, message: "Документ удален"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

export async function GetOrganisationById(id: number): Promise<Organisation | null> {
    try {
        const res = await api.get(`/org/${id}`);
        return res.data;
    } catch (e) {
        console.warn(e);
        return null;
    }
}