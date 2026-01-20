'use client'

import React, {useEffect, useState} from 'react';
import {GetMyOrganisations, Organisation} from "@/api/organisation";
import Link from "next/link";
import {Card, CardContent, CardHeader, CardTitle} from "@/app/components/ui/card";
import {Button} from "@/app/components/ui/button";
import {Plus, Building2, AlertTriangle, CheckCircle, Clock, Briefcase, UserCircle} from "lucide-react";
import styles from "@/app/user-profile/page.module.css";
import Image from "next/image";
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
            <div className="min-h-screen bg-[#989694] flex items-center justify-center">
                <Spinner />
            </div>
        );
    }

    return (
        <div style={{minHeight: '100vh', backgroundColor: '#989694', padding: '3rem 2rem'}}>
            <div className="max-w-6xl mx-auto">
                <div className="flex justify-between items-center mb-8 flex-wrap gap-4">
                    <div>
                        <h1 className="text-3xl font-bold text-white mb-2 uppercase font-montserrat">Мои организации</h1>
                        <p className="text-gray-200">Управляйте своими юридическими лицами и счетами</p>
                    </div>
                    <Link href="/organisation/create">
                        <Button className="bg-[#825e9c] text-white font-bold hover:bg-[#6a4c80] shadow-lg">
                            <Plus className="w-4 h-4 mr-2"/> Создать новую
                        </Button>
                    </Link>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {/* Карточка создания новой */}
                    <Link href="/organisation/create" className="group h-full">
                        <div className="h-full min-h-[220px] border-dashed border-2 border-gray-600 bg-transparent flex flex-col items-center justify-center rounded-lg cursor-pointer hover:border-[#825e9c] hover:bg-white/5 transition-all p-6">
                            <div className="h-16 w-16 rounded-full bg-gray-700 flex items-center justify-center mb-4 group-hover:bg-[#825e9c] text-gray-400 group-hover:text-white transition-colors">
                                <Plus size={32} />
                            </div>
                            <h3 className="text-xl font-bold text-gray-300 group-hover:text-white uppercase">Добавить организацию</h3>
                        </div>
                    </Link>

                    {/* Список организаций */}
                    {orgs.map(org => {
                        const avatarUrl = getStorageUrl(org.avatar_path, BUCKETS.AVATARS);
                        const OrgIcon = org.org_type === 'ip' ? Briefcase : org.org_type === 'jur' ? Building2 : UserCircle;

                        return (
                            <Link href={`/organisation/${org.id}`} key={org.id} className="block h-full">
                                <div className={`h-full bg-[#656662] rounded-lg shadow-lg border border-gray-500 overflow-hidden hover:translate-y-[-5px] transition-transform hover:shadow-2xl hover:border-[#DB935B] flex flex-col ${org.is_banned ? 'border-red-500' : ''}`}>

                                    {/* Header */}
                                    <div className="p-4 border-b border-gray-500 flex items-center gap-4">
                                        {/* Аватар с фиксированным размером */}
                                        <div className="relative w-12 h-12 flex-shrink-0 rounded-full overflow-hidden border-2 border-[#825e9c] bg-white">
                                            <Image
                                                src={avatarUrl || placeholder}
                                                alt={org.name}
                                                fill
                                                style={{objectFit: 'cover'}}
                                                unoptimized
                                            />
                                        </div>

                                        <div className="flex-1 min-w-0">
                                            <h3 className="text-white font-bold text-lg truncate mb-1" title={org.name}>
                                                {org.name}
                                            </h3>
                                            <div className="flex items-center gap-2">
                                                <span className="text-xs uppercase bg-gray-700 text-gray-300 px-2 py-0.5 rounded font-bold">
                                                    {org.org_type}
                                                </span>
                                                {org.is_banned && <Badge variant="destructive" className="h-5 text-[10px]">BAN</Badge>}
                                            </div>
                                        </div>
                                    </div>

                                    {/* Content */}
                                    <div className="p-4 flex-1 flex flex-col gap-3">
                                        <div className="flex justify-between items-center bg-[#555652] p-2 rounded">
                                            <span className="text-xs text-gray-400 uppercase">Баланс</span>
                                            <span className="font-bold text-[#DB935B]">{(org.balance ?? 0).toLocaleString()} ₽</span>
                                        </div>

                                        <div className="mt-auto">
                                            {org.registration_completed ? (
                                                <div className="flex items-center text-green-400 gap-2 text-sm font-medium">
                                                    <CheckCircle size={16} />
                                                    <span>Активна</span>
                                                </div>
                                            ) : (
                                                <div className="flex items-center text-yellow-500 gap-2 text-sm font-medium">
                                                    <Clock size={16} />
                                                    <span>Требует заполнения</span>
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            </Link>
                        );
                    })}
                </div>
            </div>
        </div>
    );
}