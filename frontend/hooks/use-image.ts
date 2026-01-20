// frontend/hooks/use-image.ts
import { useState, useEffect } from 'react';
import { StaticImageData } from 'next/image';
import { getStorageUrl } from "@/lib/config";

export function useImage(
    file: File | null | undefined,         // Файл, выбранный пользователем (blob)
    backendPath: string | null | undefined,// Путь, пришедший с бэкенда (строка)
    bucket: string,                        // Имя бакета
    fallback: string | StaticImageData     // Заглушка
) {
    // Определяем начальное состояние
    const getInitialSrc = () => {
        if (backendPath) return getStorageUrl(backendPath, bucket);
        return fallback;
    };

    const [src, setSrc] = useState<string | StaticImageData>(getInitialSrc());

    useEffect(() => {
        // 1. Приоритет: Локальный файл (пользователь только что выбрал)
        if (file) {
            const objectUrl = URL.createObjectURL(file);
            setSrc(objectUrl);

            // Чистим память при размонтировании или смене файла
            return () => URL.revokeObjectURL(objectUrl);
        }

        // 2. Если файла нет, но есть путь с сервера
        if (backendPath) {
            setSrc(getStorageUrl(backendPath, bucket));
            return;
        }

        // 3. Иначе заглушка
        setSrc(fallback);
    }, [file, backendPath, bucket, fallback]);

    return src;
}