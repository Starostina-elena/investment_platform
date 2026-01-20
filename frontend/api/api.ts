import axios from "axios";
import {useUserStore} from "@/context/user-store";
import {atob} from "next/dist/compiled/@edge-runtime/primitives";


// const BACK_URL = 'http://gateway:80/api';
const BACK_URL = '/api';

export interface Message {
    message: string;
    isError: boolean;
}

export const api = axios.create({
    baseURL: BACK_URL
})

export const serverApi = axios.create({
    baseURL: 'http://gateway:80/api' // TODO: Change to real server url
})

api.interceptors.request.use(function (config: any) {
    const token = useUserStore.getState().token;
    if(token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, function (error: any) {
    return Promise.reject(error);
});

export function DefaultErrorHandler(setError: (err: Message) => void) {
    return (err: any) => {
        if (err.code == 403)
            useUserStore.getState().Logout();
        const message = err.response.data;
        setError({isError: true, message: JSON.stringify(message)});
    }
}