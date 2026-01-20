import { useState, useEffect } from 'react';
import { StaticImageData } from 'next/image';
import { getStorageUrl } from "@/lib/config";

// Возвращает src для картинки
export function useImage(
    file: File | null | undefined,
    backendPath: string | null | undefined,
    bucket: string,
    fallback: string | StaticImageData
) {
    const [src, setSrc] = useState<string | StaticImageData>(fallback);

    useEffect(() => {
        if (file) {
            const objectUrl = URL.createObjectURL(file);
            setSrc(objectUrl);
            return () => URL.revokeObjectURL(objectUrl);
        }

        if (backendPath) {
            // ВАЖНО: Добавляем timestamp, чтобы сбросить кэш браузера при обновлении аватара
            // (если имя файла одинаковое, например orgpic_1.jpg)
            const url = getStorageUrl(backendPath, bucket);
            setSrc(`${url}?t=${new Date().getTime()}`);
            return;
        }

        setSrc(fallback);
    }, [file, backendPath, bucket, fallback]);

    return src;
}

// Специальный хук для проекта (у него может быть старая логика)
export function useProjectImage(
    file: File | null | undefined,
    backendPath: string | null | undefined,
    fallback: string | StaticImageData
) {
    // Просто обертка, указывающая бакет проектов
    // Но лучше использовать BUCKETS.PROJECTS явно в useImage
    return useImage(file, backendPath, 'projects', fallback);
}