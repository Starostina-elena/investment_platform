'use client'

import React, {useEffect, useState} from 'react';
import {GetMyOrganisations, Organisation} from "@/api/organisation";
import Link from "next/link";
import {Card, CardContent, CardHeader, CardTitle} from "@/app/components/ui/card";
import {Button} from "@/app/components/ui/button";
import {Plus} from "lucide-react";
import styles from "@/app/user-profile/page.module.css";
import Image from "next/image";
import {getStorageUrl, BUCKETS} from "@/lib/config";
import placeholder from "@/public/image_bg.png";

export default function MyOrganisationsPage() {
    const [orgs, setOrgs] = useState<Organisation[]>([]);

    useEffect(() => {
        GetMyOrganisations().then(setOrgs);
    }, []);

    return (
        <div className={styles.container} style={{padding: '3rem'}}>
            <div className="max-w-6xl mx-auto">
                <div className="flex justify-between items-center mb-8">
                    <h1 className="text-3xl font-bold text-white">Мои организации</h1>
                    <Link href="/organisation/create">
                        <Button className="bg-light-blue text-black font-bold">
                            <Plus className="w-4 h-4 mr-2"/> Создать новую
                        </Button>
                    </Link>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {orgs.map(org => (
                        <Link href={`/organisation/${org.id}`} key={org.id}>
                            <Card className="hover:scale-[1.02] transition-transform cursor-pointer bg-[#1e0e31] border-gray-700">
                                <CardHeader className="flex flex-row items-center gap-4">
                                    <div className="relative w-12 h-12 rounded-full overflow-hidden border border-gray-500">
                                        <Image
                                            src={getStorageUrl(org.avatar_path, BUCKETS.AVATARS) || placeholder}
                                            alt={org.name} fill className="object-cover"
                                        />
                                    </div>
                                    <div>
                                        <CardTitle className="text-white text-lg">{org.name}</CardTitle>
                                        <p className="text-sm text-gray-400">{org.org_type.toUpperCase()}</p>
                                    </div>
                                </CardHeader>
                                <CardContent>
                                    <p className="text-gray-300">Баланс: <span className="text-light-green font-bold">{org.balance} ₽</span></p>
                                    <p className="text-sm text-gray-500 mt-2">
                                        {org.registration_completed ? 'Активна' : 'Требует заполнения данных'}
                                    </p>
                                </CardContent>
                            </Card>
                        </Link>
                    ))}

                    {orgs.length === 0 && (
                        <div className="col-span-full text-center text-gray-400 py-20">
                            У вас пока нет организаций. Создайте первую, чтобы запускать проекты!
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}