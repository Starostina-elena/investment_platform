import {api, DefaultErrorHandler, Message} from "@/api/api";
import {useUserStore} from "@/context/user-store";

export async function Register(payload: any, setMessage: (message: Message) => void) {
    try {
        await api.post("/user/create", payload); // Бэкенд эндпоинт
        // После регистрации обычно нужно логиниться отдельно или бэк сразу возвращает токен.
        // В вашем коде handler/auth.go LoginHandler возвращает токен.
        // CreateUserHandler возвращает созданного юзера без токена.
        // Поэтому после регистрации делаем автоматический логин:
        await Login(payload.email, payload.password, setMessage);
    } catch (err: any) {
        DefaultErrorHandler(setMessage)(err);
    }
}

export async function Login(email: string, password: string, setMessage: (message: Message) => void) {
    try {
        const res = await api.post("/user/login", { email, password });

        const { access_token, user_id } = res.data;

        // Получаем полные данные пользователя
        const userRes = await api.get(`/user/${user_id}`, {
            headers: { Authorization: `Bearer ${access_token}` }
        });

        useUserStore.getState().Login(userRes.data, access_token);
        setMessage({isError: false, message: "Вы вошли в аккаунт!"});
    } catch (err: any) {
        DefaultErrorHandler(setMessage)(err);
    }
}

export async function UpdateUser(payload: any, setError: (err: Message) => void) {
    throw new Error("Not implemented");
    try {
        const res = await api.put('/user/update', payload);
        setError({isError: false, message: "Данные обновлены"});
        return res.data;
    } catch (err: any) {
        DefaultErrorHandler(setError)(err);
        return null;
    }
}

export async function DeleteUser(id?: number) {
    throw new Error("Not implemented");
    if (id == undefined)
        id = useUserStore.getState().user?.id;
    try {
        const res = await api.delete('/user/delete/' + id);
    } catch (err: any) {
        DefaultErrorHandler(() => {})(err);
    }
}
