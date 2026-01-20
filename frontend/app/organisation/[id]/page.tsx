'use client'

import React, {useEffect, useRef, useState} from 'react';
import {useParams} from "next/navigation";
import {GetFullOrganisation, Organisation, UpdateOrganisation, UploadOrgAvatar} from "@/api/organisation";
import OrganisationForm from "@/app/components/organisation-form";
import Spinner from "@/app/components/spinner";
import styles from "@/app/user-profile/page.module.css";
import {useImage} from "@/hooks/use-image";
import {BUCKETS} from "@/lib/config";
import Image from "next/image";
import {Button} from "@/app/components/ui/button";
import {Camera, Edit3} from "lucide-react";
import MessageComponent from "@/app/components/message";
import {Message} from "@/api/api";

// Заглушка для аватара орги
import orgPlaceholder from "@/public/image_bg.png";

export default function OrganisationPage() {
    const params = useParams();
    const [org, setOrg] = useState<Organisation | null | undefined>(undefined);
    const [isEditing, setIsEditing] = useState(false);
    const [message, setMessage] = useState<Message | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    // Загрузка данных
    useEffect(() => {
        if (params.id) {
            GetFullOrganisation(+params.id).then(setOrg);
        }
    }, [params.id]);

    const avatarSrc = useImage(null, org?.avatar_path, BUCKETS.AVATARS, orgPlaceholder);

    const handleUpdate = async (data: any) => {
        if (!org) return;
        const updated = await UpdateOrganisation(org.id, data, setMessage);
        if (updated) {
            setOrg({...org, ...updated}); // Merge updates
            setIsEditing(false);
        }
    };

    const handleAvatarUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file && org) {
            const path = await UploadOrgAvatar(org.id, file, setMessage);
            if (path) {
                setOrg({...org, avatar_path: path});
            }
        }
    };

    if (org === undefined) return <div className={styles.container}><Spinner/></div>;
    if (org === null) return <div className={styles.container}><h1 className="text-white text-center pt-20">Организация не найдена или нет доступа</h1></div>;

    return (
        <div className={styles.container} style={{padding: '2rem'}}>
            <div className="max-w-5xl mx-auto grid grid-cols-1 lg:grid-cols-3 gap-8">

                {/* Левая колонка - Аватар и статус */}
                <div className="lg:col-span-1">
                    <div className="bg-[#1e0e31] p-6 rounded-lg shadow-lg text-center">
                        <div className="relative w-40 h-40 mx-auto mb-4 rounded-full overflow-hidden border-4 border-light-green">
                            <Image src={avatarSrc} alt={org.name} fill className="object-cover"/>
                        </div>

                        <input type="file" ref={fileInputRef} hidden onChange={handleAvatarUpload} accept="image/*"/>
                        <Button variant="outline" className="mb-4" onClick={() => fileInputRef.current?.click()}>
                            <Camera className="w-4 h-4 mr-2"/> Сменить лого
                        </Button>

                        <h2 className="text-xl font-bold text-white mb-2">{org.name}</h2>
                        <div className={`inline-block px-3 py-1 rounded text-sm font-bold ${org.registration_completed ? 'bg-green-500 text-black' : 'bg-yellow-500 text-black'}`}>
                            {org.registration_completed ? 'Активна' : 'На модерации'}
                        </div>

                        <div className="mt-6 text-left text-gray-300 space-y-2 text-sm">
                            <p><strong>Баланс:</strong> {org.balance} ₽</p>
                            <p><strong>Тип:</strong> {org.org_type.toUpperCase()}</p>
                            <p><strong>Создана:</strong> {org.created_at ? new Date(org.created_at).toLocaleDateString() : '-'}</p>
                        </div>
                    </div>
                </div>

                {/* Правая колонка - Данные */}
                <div className="lg:col-span-2">
                    <div className="flex justify-between items-center mb-6">
                        <h1 className="text-3xl font-bold text-white">Профиль организации</h1>
                        {!isEditing && (
                            <Button onClick={() => setIsEditing(true)}>
                                <Edit3 className="w-4 h-4 mr-2"/> Редактировать
                            </Button>
                        )}
                        {isEditing && (
                            <Button variant="ghost" className="text-white" onClick={() => setIsEditing(false)}>
                                Отмена
                            </Button>
                        )}
                    </div>

                    <div className="bg-[#1e0e31] p-1 rounded-lg">
                        <OrganisationForm
                            initialData={org}
                            onSubmit={handleUpdate}
                            isEditing={isEditing}
                        />
                        {/* При просмотре блокируем форму CSS-ом или пропсом (в OrganisationForm можно добавить disabled={!isEditing}) */}
                        {!isEditing && <div className="absolute inset-0 z-10 bg-transparent" />}
                    </div>

                    <div className="mt-4">
                        <MessageComponent message={message}/>
                    </div>
                </div>
            </div>
        </div>
    );
}