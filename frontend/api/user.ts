import {api, DefaultErrorHandler, Message} from "@/api/api";

export interface User {
    id: number;
    name: string;
    surname: string;
    patronymic?: string;
    nickname: string;
    email: string;
    balance: number;
    avatar_path?: string;
    is_admin: boolean;
    is_banned: boolean;
    created_at: string;
}

export async function GetUserById(id: number): Promise<User | null> {
    try {
        const res = await api.get(`/user/${id}`);
        const user = res.data;
        // Хак: Бэкенд не отдает поле avatar_path в JSON.
        // Формируем стандартный путь, если аватар есть (проверка на существование будет в компоненте через onError или просто попытку загрузки)
        // В идеале бэкенд должен отдавать это поле.
        if (!user.avatar_path) {
            user.avatar_path = `userpic_${user.id}.jpg`;
        }
        return user;
    } catch (e) {
        console.warn("Failed to fetch user:", e);
        return null;
    }
}

export async function UpdateUserInfo(user: Partial<User>, setMessage: (msg: Message) => void): Promise<User | null> {
    try {
        // Отправляем только разрешенные поля
        const payload = {
            name: user.name,
            surname: user.surname,
            patronymic: user.patronymic,
            nickname: user.nickname,
            email: user.email
        };

        const res = await api.post('/user/update', payload);
        setMessage({isError: false, message: "Профиль обновлен"});
        return res.data;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return null;
    }
}

export async function UploadUserAvatar(file: File, setMessage: (msg: Message) => void): Promise<string | null> {
    try {
        const formData = new FormData();
        formData.append("avatar", file);

        const res = await api.post('/user/avatar/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });

        setMessage({isError: false, message: "Аватар обновлен"});
        return res.data.avatar_path;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return null;
    }
}

export async function BanUser(userId: number, ban: boolean, setMessage: (msg: Message) => void): Promise<boolean> {
    try {
        await api.post(`/user/${userId}/active?ban=${ban}`);
        setMessage({isError: false, message: ban ? "Пользователь заблокирован" : "Пользователь разблокирован"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

export interface Investment {
    project_id: number;
    project_name: string;
    quick_peek: string;
    monetization_type: string;
    total_invested: number;
    total_received: number;
    created_at: string;
    is_completed: boolean;
    is_banned: boolean;
}

export async function GetActiveInvestments(): Promise<Investment[]> {
    try {
        const res = await api.get('/user/investments/active');
        return Array.isArray(res.data) ? res.data : [];
    } catch (e) {
        console.warn("Failed to fetch active investments:", e);
        return [];
    }
}

export async function GetArchivedInvestments(): Promise<Investment[]> {
    try {
        const res = await api.get('/user/investments/archived');
        return Array.isArray(res.data) ? res.data : [];
    } catch (e) {
        console.warn("Failed to fetch archived investments:", e);
        return [];
    }
}