'use client'

import React, {useEffect, useState} from 'react';
import {GetMyOrganisations, Organisation} from "@/api/organisation";
import Link from "next/link";
import {Card, CardContent, CardHeader, CardTitle} from "@/app/components/ui/card";
import {Button} from "@/app/components/ui/button";
import {Plus, Building2, AlertTriangle, CheckCircle, Clock} from "lucide-react";
import styles from "@/app/user-profile/page.module.css"; // Переиспользуем стили профиля
import {getStorageUrl, BUCKETS} from "@/lib/config";
import placeholder from "@/public/image_bg.png";
import Spinner from "@/app/components/spinner";
import {Badge} from "@/app/components/ui/badge";

export default function MyOrganisationsPage() {
    const [orgs, setOrgs] = useState<Organisation[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        GetMyOrganisations()
            .then(setOrgs)
            .finally(() => setLoading(false));
    }, []);

    if (loading) {
        return (
            <div className="min-h-screen bg-[#130622] flex items-center justify-center">
                <Spinner />
            </div>
        );
    }

    return (
        <div className={styles.container} style={{padding: '3rem', minHeight: '100vh'}}>
            <div className="max-w-6xl mx-auto">
                <div className="flex justify-between items-center mb-8">
                    <div>
                        <h1 className="text-3xl font-bold text-white mb-2">Мои организации</h1>
                        <p className="text-gray-400">Управляйте своими юридическими лицами и счетами</p>
                    </div>
                    <Link href="/organisation/create">
                        <Button className="bg-[#825e9c] text-black font-bold hover:bg-[#00b0d6]">
                            <Plus className="w-4 h-4 mr-2"/> Создать новую
                        </Button>
                    </Link>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {orgs.map(org => (
                        <Link href={`/organisation/${org.id}`} key={org.id} className="group">
                            <Card className={`h-full border-gray-700 bg-[#1e0e31] transition-all duration-300 group-hover:border-light-blue/50 group-hover:-translate-y-1 ${org.is_banned ? 'border-red-500/50' : ''}`}>
                                <CardHeader className="flex flex-row items-start gap-4 pb-2">
                                    <div className="relative w-14 h-14 shrink-0 rounded-full overflow-hidden border-2 border-gray-600 group-hover:border-light-blue">
                                        <img
                                            src={getStorageUrl(org.avatar_path, BUCKETS.AVATARS) || placeholder.src}
                                            alt={org.name}
                                            style={{width: '100%', height: '100%', objectFit: 'cover', pointerEvents: 'none', userSelect: 'none'}}
                                            draggable={false}
                                        />
                                    </div>
                                    <div className="flex-1 min-w-0">
                                        <div className="flex items-start justify-between gap-2">
                                            <CardTitle className="text-white text-lg truncate pr-2 group-hover:text-light-blue transition-colors">
                                                {org.name}
                                            </CardTitle>
                                            {org.is_banned && (
                                                <Badge variant="destructive" className="shrink-0 h-6">BAN</Badge>
                                            )}
                                        </div>
                                        <p className="text-xs uppercase font-bold tracking-wider text-gray-500 mt-1">
                                            {org.org_type === 'ip' ? 'ИП' : org.org_type === 'jur' ? 'Юр. лицо' : 'Физ. лицо'}
                                        </p>
                                    </div>
                                </CardHeader>
                                <CardContent>
                                    <div className="space-y-3">
                                        <div className="flex items-center justify-between p-2 rounded bg-white/5">
                                            <span className="text-sm text-gray-400">Баланс</span>
                                            <span className="font-bold text-[#DB935B]">{(org.balance ?? 0).toLocaleString()} ₽</span>
                                        </div>

                                        <div className="flex items-center gap-2 text-sm">
                                            {org.registration_completed ? (
                                                <div className="flex items-center text-green-400 gap-2">
                                                    <CheckCircle size={16} />
                                                    <span>Активна</span>
                                                </div>
                                            ) : (
                                                <div className="flex items-center text-yellow-500 gap-2">
                                                    <Clock size={16} />
                                                    <span>Не заполнена</span>
                                                </div>
                                            )}
                                        </div>

                                        {!org.registration_completed && (
                                            <div className="text-xs text-yellow-500/80 bg-yellow-500/10 p-2 rounded border border-yellow-500/20">
                                                Заполните реквизиты, чтобы запускать проекты.
                                            </div>
                                        )}

                                        {org.is_banned && (
                                            <div className="flex items-center gap-2 text-red-400 text-xs bg-red-500/10 p-2 rounded border border-red-500/20">
                                                <AlertTriangle size={14} />
                                                <span>Организация заблокирована администрацией.</span>
                                            </div>
                                        )}
                                    </div>
                                </CardContent>
                            </Card>
                        </Link>
                    ))}

                    {/* Карточка создания новой (пустая) */}
                    <Link href="/organisation/create" className="group">
                        <Card className="h-full border-dashed border-2 border-gray-700 bg-transparent flex flex-col items-center justify-center min-h-[200px] cursor-pointer hover:border-light-blue hover:bg-white/5 transition-all">
                            <div className="h-16 w-16 rounded-full bg-gray-800 flex items-center justify-center mb-4 group-hover:bg-[#825e9c] group-hover:text-black transition-colors text-gray-400">
                                <Plus size={32} />
                            </div>
                            <h3 className="text-xl font-bold text-gray-400 group-hover:text-white">Добавить организацию</h3>
                        </Card>
                    </Link>
                </div>
            </div>
        </div>
    );
}