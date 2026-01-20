'use client'
import { Card, CardContent, CardHeader, CardTitle } from "@/app/components/ui/card";
import { Badge } from "@/app/components/ui/badge";
import { Investment } from "@/api/user";
import Link from "next/link";
import { ArrowRight, TrendingUp, CheckCircle, Lock } from "lucide-react";

export default function InvestmentsList({ investments }: { investments: Investment[] }) {
    if (investments.length === 0) {
        return <div className="text-center text-gray-500 py-8">Инвестиций не найдено</div>;
    }

    return (
        <div className="grid gap-4">
            {investments.map((inv) => (
                <Card key={inv.project_id} className="bg-[#656662] border-gray-700 hover:border-gray-500 transition-colors">
                    <CardContent className="p-4 flex flex-col md:flex-row gap-4 justify-between items-start md:items-center">
                        <div className="flex-1">
                            <div className="flex items-center gap-2 mb-2">
                                <Link href={`/project?id=${inv.project_id}`} className="text-lg font-bold text-white hover:text-[#825e9c] hover:underline">
                                    {inv.project_name}
                                </Link>
                                {inv.is_banned && <Badge variant="destructive" className="text-xs">Заблокирован</Badge>}
                                {inv.is_completed && <Badge variant="secondary" className="text-xs bg-green-900 text-green-100">Завершен</Badge>}
                            </div>
                            <p className="text-sm text-gray-400 mb-2">{inv.quick_peek}</p>
                            <div className="flex gap-2 text-xs">
                                <span className="bg-white/10 px-2 py-1 rounded text-gray-300">
                                    {formatMonetization(inv.monetization_type)}
                                </span>
                                <span className="text-gray-500 py-1">
                                    {new Date(inv.created_at).toLocaleDateString()}
                                </span>
                            </div>
                        </div>

                        <div className="flex flex-col items-end gap-1 min-w-[150px]">
                            <div className="text-sm text-gray-400">Вложено:</div>
                            <div className="text-xl font-bold text-white">{inv.total_invested} ₽</div>

                            {inv.monetization_type !== 'charity' && (
                                <>
                                    <div className="text-sm text-gray-400 mt-1">Получено:</div>
                                    <div className={`text-lg font-bold ${inv.total_received > 0 ? 'text-green-400' : 'text-gray-500'}`}>
                                        {inv.total_received} ₽
                                    </div>
                                </>
                            )}
                        </div>
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}

function formatMonetization(type: string): string {
    switch (type) {
        case 'charity': return 'Благотворительность';
        case 'fixed_percent': return 'Фикс. процент';
        case 'time_percent': return 'Временной процент';
        case 'custom': return 'Свое вознаграждение';
        default: return type;
    }
}