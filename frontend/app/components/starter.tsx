'use client'
import {useEffect} from "react";
import {useUserStore} from "@/context/user-store";
import {usePathname, useRouter} from "next/navigation";

export default function Starter(){
    const router = useRouter();
    const path = usePathname();
    const init = useUserStore((state) => state.init);

    useEffect(() => {
        // Вызываем инициализацию, которая проверит localStorage
        // и восстановит сессию, либо редиректнет на логин, если нужно
        init(router, path);
    }, [path, router, init]);

    return <></>
}