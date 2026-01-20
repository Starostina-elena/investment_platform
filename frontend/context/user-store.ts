import {create} from "zustand";
import {AppRouterInstance} from "next/dist/shared/lib/app-router-context.shared-runtime";
import {addBasePath} from "next/dist/client/add-base-path";
import {User} from "@/api/user";

const AllowedRoutes = ['/login', '/', ''];

interface UserState {
    inited: boolean;
    user: User | null;
    token: string | null;
    Login: (user: User, token: string) => void;
    Logout: () => void;
    setToken: (token: string) => void; // <--- Добавили
    init: (router: AppRouterInstance, pathname: string) => void;
}

export const useUserStore = create<UserState>((set) => ({
    inited: false,
    user: null,
    token: null,
    Login: function (user, token) {
        localStorage.setItem("user", JSON.stringify(user));
        localStorage.setItem("token", token);
        set({user, token: token});
    },
    Logout: function () {
        localStorage.removeItem("user");
        localStorage.removeItem("token");
        // Также можно дернуть API logout, чтобы очистить cookie на сервере
        // api.post('/user/logout').catch(() => {});
        set({user: null, token: null});
    },
    setToken: function(token) { // <--- Реализация
        localStorage.setItem("token", token);
        set({token: token});
    },
    init: function(router, pathname) {
        try {
            const userStr = localStorage.getItem("user");
            const token = localStorage.getItem("token");

            if (userStr && token) {
                const user: User = JSON.parse(userStr);
                set({token: token, inited: true, user});
            } else {
                set({inited: true});
                // Проверка маршрутов: если страница защищенная, редиректим
                // Но лучше это делать в middleware Next.js или в HOC компонентах
                // if(!AllowedRoutes.map(e => addBasePath(e)).includes(pathname.replace(/\/$/, ''))) {
                //     router.push("/login");
                // }
            }
        } catch (e) {
            console.warn(e);
            set({inited: true});
            // router.push("/login");
        }
    }
}));