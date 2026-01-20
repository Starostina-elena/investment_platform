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