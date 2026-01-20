// frontend/api/project.ts
import {api, DefaultErrorHandler, Message} from "@/api/api";

// Соответствует services/project/core/models.go
export interface Project {
    id: number;
    name: string;
    creator_id: number; // ID организации
    quick_peek: string;
    quick_peek_picture_path?: string;
    content: string;
    is_public: boolean;
    is_completed: boolean;
    current_money: number;
    wanted_money: number;
    duration_days: number;
    created_at?: string;
    is_banned: boolean;
    monetization_type: string;
    percent?: number;

    // UI поля (не отправляются в JSON проекта напрямую, но используются для логики)
    quickPeekPictureFile?: File | null;
    category?: string; // Пока бэк не хранит, но фронт использует
    location?: string; // Пока бэк не хранит
}

export async function GetProjectById(id: number): Promise<Project | null> {
    try {
        const res = await api.get(`/projects/${id}`); // Nginx: /api/projects/ -> project_service
        return res.data;
    } catch (e) {
        console.warn(e);
        return null;
    }
}

export async function GetProjects(
    limit: number = 50,
    offset: number = 0,
    query?: string,
    category?: string
): Promise<Project[]> {
    try {
        const params = new URLSearchParams();
        params.append("limit", limit.toString());
        params.append("offset", offset.toString());
        if (query) params.append("q", query);
        if (category) params.append("category", category);

        const res = await api.get(`/projects/?${params.toString()}`);
        if (Array.isArray(res.data)) {
            return res.data;
        }
        return [];
    } catch (e) {
        console.warn(e);
        return [];
    }
}

export async function PublishProject(project: Project, setMessage: (message: Message) => void) {
    try {
        // 1. Формируем payload строго по Go struct CreateProjectRequest
        const payload = {
            name: project.name,
            creator_id: project.creator_id,
            quick_peek: project.quick_peek,
            content: project.content,
            wanted_money: project.wanted_money,
            duration_days: project.duration_days,
            monetization_type: project.monetization_type || "donation",
            percent: project.percent || 0
        };

        const res = await api.post('/projects/create', payload);
        const createdProject = res.data;

        // 2. Если пользователь выбрал картинку, грузим её отдельным запросом
        if (project.quickPeekPictureFile) {
            const formData = new FormData();
            formData.append("picture", project.quickPeekPictureFile);

            // Go handler: UploadPictureHandler ожидает "picture" в form-data
            await api.post(`/projects/${createdProject.id}/picture/upload`, formData, {
                headers: { 'Content-Type': 'multipart/form-data' }
            });
        }

        setMessage({isError: false, message: "Проект успешно создан!"});
        return createdProject.id;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return null;
    }
}

export async function BanProject(projectId: number, ban: boolean, setMessage: (msg: Message) => void): Promise<boolean> {
    try {
        await api.post(`/projects/${projectId}/ban?ban=${ban}`);
        setMessage({isError: false, message: ban ? "Проект заблокирован" : "Проект разблокирован"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

