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
    monetization_type: string;
    created_at?: string;
    is_banned: boolean;
    percent?: number;

    quickPeekPictureFile?: File | null;

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
        if (query) params.append("search", query);
        if (category) params.append("type", category); // Используем type вместо category

        const res = await api.get(`/projects/projects?${params.toString()}`);
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
        // При отправке используем snake_case
        const payload = {
            creator_id: project.creator_id,
            name: project.name,
            quick_peek: project.quick_peek,
            content: project.content,
            wanted_money: project.wanted_money,
            duration_days: project.duration_days,
            monetization_type: project.monetization_type || "charity",
            percent: project.percent || 0
        };

        const res = await api.post('/projects/create', payload);
        const createdProject = res.data;

        if (project.quickPeekPictureFile) {
            const formData = new FormData();
            formData.append("picture", project.quickPeekPictureFile);

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

export async function UpdateProject(
    projectId: number,
    project: Partial<Project>,
    setMessage: (msg: Message) => void
): Promise<Project | null> {
    try {
        const res = await api.post(`/projects/${projectId}/update`, project);
        setMessage({isError: false, message: "Проект обновлен"});
        return res.data;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return null;
    }
}

export async function SetProjectCompleted(
    projectId: number,
    completed: boolean,
    setMessage: (msg: Message) => void
): Promise<boolean> {
    try {
        await api.post(`/projects/${projectId}/completed?completed=${completed}`);
        setMessage({isError: false, message: completed ? "Проект завершен" : "Проект возобновлен"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

export async function SetProjectPublic(
    projectId: number,
    isPublic: boolean,
    setMessage: (msg: Message) => void
): Promise<boolean> {
    try {
        await api.post(`/projects/${projectId}/public?public=${isPublic}`);
        setMessage({isError: false, message: isPublic ? "Проект опубликован" : "Проект скрыт"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

export async function ProjectPayback(
    projectId: number,
    setMessage: (msg: Message) => void
): Promise<boolean> {
    try {
        await api.post(`/projects/${projectId}/payback`);
        setMessage({isError: false, message: "Выплата дивидендов выполнена"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

export async function GetOrganisationProjects(orgId: number): Promise<Project[]> {
    try {
        const res = await api.get(`/projects/projects/org/${orgId}`);
        return Array.isArray(res.data) ? res.data : [];
    } catch (e) {
        console.warn(e);
        return [];
    }
}

export async function GetAllOrganisationProjects(orgId: number): Promise<Project[]> {
    try {
        const res = await api.get(`/projects/projects/all/org/${orgId}`);
        return Array.isArray(res.data) ? res.data : [];
    } catch (e) {
        console.warn(e);
        return [];
    }
}

export async function DeleteProjectPicture(
    projectId: number,
    setMessage: (msg: Message) => void
): Promise<boolean> {
    try {
        await api.delete(`/projects/${projectId}/picture`);
        setMessage({isError: false, message: "Обложка удалена"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

export async function UploadProjectPicture(
    projectId: number,
    file: File,
    setMessage: (msg: Message) => void
): Promise<string | null> {
    try {
        const formData = new FormData();
        formData.append("picture", file);
        
        const res = await api.post(`/projects/${projectId}/picture/upload`, formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
        
        setMessage({isError: false, message: "Обложка обновлена"});
        
        // Бэкенд возвращает строку пути или объект с path
        if (typeof res.data === 'string') return res.data;
        if (res.data?.path && typeof res.data.path === 'string') return res.data.path;
        
        return null;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return null;
    }
}
