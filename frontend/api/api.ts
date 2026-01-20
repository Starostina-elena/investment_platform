import axios from "axios";
import {useUserStore} from "@/context/user-store";

// const BACK_URL = 'http://gateway:80/api';
const BACK_URL = '/api';

export interface Message {
    message: string;
    isError: boolean;
}

export const api = axios.create({
    baseURL: BACK_URL,
    withCredentials: true // Важно для отправки cookies (refresh token)
})

export const serverApi = axios.create({
    baseURL: 'http://gateway:80/api'
})

// Request Interceptor: Добавляем токен
api.interceptors.request.use(function (config: any) {
    const token = useUserStore.getState().token;
    if(token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, function (error: any) {
    return Promise.reject(error);
});

// Response Interceptor: Ловим 401 и обновляем токен
api.interceptors.response.use((response) => {
    return response;
}, async (error) => {
    const originalRequest = error.config;

    // Если ошибка 401 и мы еще не пробовали обновить токен для этого запроса
    if (error.response?.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true;

        try {
            // Пытаемся обновить токен
            // Важно: этот запрос идет с cookie, которую браузер сам подставит
            const rs = await axios.post(`${BACK_URL}/user/refresh`, {}, { withCredentials: true });

            const { access_token } = rs.data;

            // Обновляем токен в сторе (но не перезаписываем юзера целиком, чтобы не было мерцаний)
            // В useUserStore добавим метод setToken
            useUserStore.getState().setToken(access_token);

            // Обновляем заголовок в упавшем запросе и повторяем его
            originalRequest.headers['Authorization'] = `Bearer ${access_token}`;
            return api(originalRequest);
        } catch (_error) {
            // Если обновить не удалось (refresh протух), разлогиниваемся
            useUserStore.getState().Logout();
            // Можно сделать редирект на логин
            // window.location.href = '/login';
            return Promise.reject(_error);
        }
    }
    return Promise.reject(error);
});

export function DefaultErrorHandler(setError: (err: Message) => void) {
    return (err: any) => {
        // 403 ловится отдельно, если нужно
        if (err.response?.status === 403) {
            // setError({isError: true, message: "Доступ запрещен"});
        }

        const message = err.response?.data || "Неизвестная ошибка";
        // Бэкенд иногда возвращает строку, иногда JSON
        const msgText = typeof message === 'string' ? message : JSON.stringify(message);

        setError({isError: true, message: msgText});
    }
}