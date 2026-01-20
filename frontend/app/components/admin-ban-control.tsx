'use client'
import { useState } from "react";
import { Button } from "@/app/components/ui/button";
import { ShieldAlert, ShieldCheck } from "lucide-react";
import { useUserStore } from "@/context/user-store";
import MessageComponent from "@/app/components/message";
import { Message } from "@/api/api";
import Spinner from "@/app/components/spinner";
import { BanUser } from "@/api/user";
import { BanOrganisation } from "@/api/organisation";
import { BanProject } from "@/api/project";

interface Props {
    entityId: number;
    entityType: 'user' | 'org' | 'project';
    isBanned: boolean;
    onUpdate: (newStatus: boolean) => void;
}

export default function AdminBanControl({ entityId, entityType, isBanned, onUpdate }: Props) {
    const { user } = useUserStore();
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState<Message | null>(null);

    // Если не админ, ничего не рендерим
    if (!user?.is_admin) return null;

    const handleToggleBan = async () => {
        setLoading(true);
        setMessage(null);
        let success = false;
        const newStatus = !isBanned;

        if (entityType === 'user') {
            success = await BanUser(entityId, newStatus, setMessage);
        } else if (entityType === 'org') {
            success = await BanOrganisation(entityId, newStatus, setMessage);
        } else if (entityType === 'project') {
            success = await BanProject(entityId, newStatus, setMessage);
        }

        if (success) {
            onUpdate(newStatus);
        }
        setLoading(false);
    };

    return (
        <div className="mt-6 border-t border-gray-700 pt-4">
            <p className="text-xs text-gray-500 mb-2 uppercase font-bold tracking-wider">Панель администратора</p>
            <div className="flex flex-col gap-2">
                <Button
                    variant={isBanned ? "default" : "destructive"}
                    onClick={handleToggleBan}
                    disabled={loading}
                    className="w-full gap-2"
                >
                    {loading ? <Spinner size={20} /> : (isBanned ? <ShieldCheck /> : <ShieldAlert />)}
                    {isBanned ? "Разблокировать" : "Заблокировать"}
                </Button>
                <MessageComponent message={message} />
            </div>
        </div>
    );
}