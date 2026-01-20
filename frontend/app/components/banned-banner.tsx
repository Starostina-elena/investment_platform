'use client'
import { AlertTriangle } from "lucide-react";

export default function BannedBanner({ type }: { type: 'Пользователь' | 'Проект' | 'Организация' }) {
    return (
        <div className="w-full bg-red-500/20 border border-red-500 rounded-lg p-4 mb-6 flex items-center gap-4 text-red-100">
            <AlertTriangle className="h-8 w-8 text-red-500" />
            <div>
                <h3 className="font-bold text-lg text-red-500">ЗАБЛОКИРОВАНО</h3>
                <p>{type} нарушил правила платформы и был заблокирован администрацией. Действия ограничены.</p>
            </div>
        </div>
    );
}