import {api, DefaultErrorHandler, Message} from "@/api/api";

export interface Comment {
    id: number;
    user_id: number;
    username: string;
    project_id: number;
    body: string;
    created_at: string;
}

export async function GetProjectComments(projectId: number, limit: number = 50, offset: number = 0): Promise<Comment[]> {
    try {
        const res = await api.get(`/comments/read/all/${projectId}?limit=${limit}&offset=${offset}`);
        return Array.isArray(res.data) ? res.data : [];
    } catch (e) {
        console.warn("Failed to load comments:", e);
        return [];
    }
}

export async function AddComment(projectId: number, body: string, setMessage: (msg: Message) => void): Promise<Comment | null> {
    try {
        const res = await api.post(`/comments/add/${projectId}`, { body });
        return res.data;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return null;
    }
}

export async function DeleteComment(commentId: number, setMessage: (msg: Message) => void): Promise<boolean> {
    try {
        await api.delete(`/comments/delete/${commentId}`);
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}