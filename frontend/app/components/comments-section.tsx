'use client'
import {useEffect, useState} from "react";
import {AddComment, Comment, DeleteComment, GetProjectComments} from "@/api/comments";
import {useUserStore} from "@/context/user-store";
import styles from "./comments-section.module.css";
import Link from "next/link";
import Spinner from "@/app/components/spinner";
import MessageComponent from "@/app/components/message";
import {Message} from "@/api/api";

export default function CommentsSection({projectId}: { projectId: number }) {
    const {user} = useUserStore();
    const [comments, setComments] = useState<Comment[]>([]);
    const [loading, setLoading] = useState(true);
    const [newComment, setNewComment] = useState("");
    const [message, setMessage] = useState<Message | null>(null);

    const loadComments = () => {
        GetProjectComments(projectId).then(data => {
            setComments(data);
            setLoading(false);
        });
    };

    useEffect(() => {
        loadComments();
    }, [projectId]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!newComment.trim()) return;

        const res = await AddComment(projectId, newComment, setMessage);
        if (res) {
            setNewComment("");
            setComments([res, ...comments]); // Добавляем новый коммент в начало
            setMessage(null);
        }
    };

    const handleDelete = async (commentId: number) => {
        if (confirm("Вы уверены, что хотите удалить этот комментарий?")) {
            const success = await DeleteComment(commentId, setMessage);
            if (success) {
                setComments(comments.filter(c => c.id !== commentId));
            }
        }
    };

    if (loading) return <Spinner/>;

    return (
        <div className={styles.container}>

            {/* Форма отправки */}
            {user ? (
                <form className={styles.form} onSubmit={handleSubmit}>
                    <textarea
                        className={styles.textarea}
                        placeholder="Напишите комментарий..."
                        value={newComment}
                        onChange={e => setNewComment(e.target.value)}
                    />
                    <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center'}}>
                        <MessageComponent message={message}/>
                        <button type="submit" className={styles.submit_btn}>Отправить</button>
                    </div>
                </form>
            ) : (
                <div className={styles.login_placeholder}>
                    <Link href="/login" className={styles.login_link}>Войдите</Link>, чтобы оставить комментарий
                </div>
            )}

            {/* Список комментариев */}
            <div className={styles.list}>
                {comments.length === 0 ? (
                    <div className={styles.empty_state}>Пока никто не оставил комментариев. Будьте первым!</div>
                ) : (
                    comments.map(c => (
                        <div key={c.id} className={styles.comment}>
                            <div className={styles.avatar_placeholder}>
                                {c.username.charAt(0).toUpperCase()}
                            </div>
                            <div className={styles.content}>
                                <div className={styles.header}>
                                    <div>
                                        <span className={styles.author}>{c.username}</span>
                                        <span className={styles.date}>
                                            {new Date(c.created_at).toLocaleDateString()} {new Date(c.created_at).toLocaleTimeString().slice(0, 5)}
                                        </span>
                                    </div>
                                    {(user && (user.id === c.user_id || user.is_admin)) && (
                                        <button
                                            className={styles.delete_btn}
                                            onClick={() => handleDelete(c.id)}
                                        >
                                            Удалить
                                        </button>
                                    )}
                                </div>
                                <div className={styles.body}>{c.body}</div>
                            </div>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
}