// frontend/lib/config.ts

// Если мы на клиенте (браузер), то путь относительный /media/...
// Если на сервере (SSR), то путь должен быть http://gateway:80/media/... (но Image компонент это сам не сделает)
export const MEDIA_URL = '/media';

export const BUCKETS = {
    PROJECTS: 'projects',
    AVATARS: 'avatars',
    DOCUMENTS: 'documents'
};

export function getStorageUrl(path?: string, bucket: string = BUCKETS.PROJECTS): string {
    if (!path) return "";

    // Если это уже полный URL (http...), возвращаем как есть
    if (path.startsWith("http") || path.startsWith("blob:") || path.startsWith("data:")) return path;

    // Убираем слэш в начале
    const cleanPath = path.startsWith('/') ? path.slice(1) : path;

    // Возвращаем путь, который Nginx перехватит: /media/bucket/file.jpg
    return `${MEDIA_URL}/${bucket}/${cleanPath}`;
}