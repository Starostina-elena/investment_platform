// frontend/lib/config.ts

// Базовый URL API (проксируется через Next.js)
export const API_URL = '/api';

// Базовый URL для файлов (проксируется Nginx на MinIO)
// Nginx конфиг: location /media/ -> proxy_pass http://minio_api/;
export const MEDIA_URL = '/media';

export const BUCKETS = {
    PROJECTS: 'projects',
    AVATARS: 'avatars',
    DOCUMENTS: 'documents'
};

/**
 * Формирует ссылку на файл.
 * @param path Имя файла из БД (например, "projpic_1.jpg")
 * @param bucket Имя бакета
 */
export function getStorageUrl(path?: string, bucket: string = BUCKETS.PROJECTS): string {
    if (!path) return "";
    if (path.startsWith("http") || path.startsWith("blob:")) return path;

    // Убираем слэш в начале пути, если он есть, чтобы не дублировать
    const cleanPath = path.startsWith('/') ? path.slice(1) : path;

    return `${MEDIA_URL}/${bucket}/${cleanPath}`;
}