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
import {Camera, Edit3, Building2, UserCircle, Briefcase} from "lucide-react";
import MessageComponent from "@/app/components/message";
import {Message} from "@/api/api";
import {Badge} from "@/app/components/ui/badge";

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

    if (org === undefined) return <div className="min-h-screen bg-[#989694] flex items-center justify-center"><Spinner/></div>;
    if (org === null) return (
        <div className="min-h-screen bg-[#989694] flex flex-col items-center justify-center text-white">
            <h1 className="text-3xl font-bold mb-4">Организация не найдена или нет доступа</h1>
            <p className="text-gray-200">Возможно, вы пытаетесь просмотреть чужую организацию без прав администратора.</p>
        </div>
    );

    // Определяем иконку типа
    const OrgIcon = org.org_type === 'ip' ? Briefcase : org.org_type === 'jur' ? Building2 : UserCircle;
    const orgTypeName = org.org_type === 'ip' ? 'Индивидуальный предприниматель' : org.org_type === 'jur' ? 'Юридическое лицо' : 'Физическое лицо';

    return (
        <div className={styles.container} style={{padding: '3rem 2rem'}}>
            <div className="max-w-6xl mx-auto grid grid-cols-1 lg:grid-cols-3 gap-8">

                {/* Левая колонка - Карточка организации */}
                <div className="lg:col-span-1">
                    <div className="bg-[#656662] p-6 rounded-lg shadow-xl border border-gray-500 text-center sticky top-24">
                        {/* Аватар */}
                        <div className="relative w-40 h-40 mx-auto mb-6 rounded-full overflow-hidden border-4 border-[#825e9c] shadow-lg bg-white">
                            <Image
                                src={avatarSrc}
                                alt={org.name}
                                fill
                                className="object-cover"
                            />
                        </div>

                        {/* Загрузка аватара */}
                        {!isEditing && (
                            <>
                                <input type="file" ref={fileInputRef} hidden onChange={handleAvatarUpload} accept="image/*"/>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    className="mb-6 border-gray-400 text-white hover:bg-white hover:text-black"
                                    onClick={() => fileInputRef.current?.click()}
                                >
                                    <Camera className="w-4 h-4 mr-2"/> Сменить логотип
                                </Button>
                            </>
                        )}

                        <h2 className="text-2xl font-bold text-white mb-2 leading-tight">{org.name}</h2>

                        <div className="flex justify-center gap-2 mb-6">
                            <Badge className={`${org.registration_completed ? 'bg-green-600' : 'bg-yellow-600'} text-white border-none`}>
                                {org.registration_completed ? 'Активна' : 'Черновик'}
                            </Badge>
                            {org.is_banned && <Badge variant="destructive">Заблокирована</Badge>}
                        </div>

                        <div className="text-left bg-[#555652] p-4 rounded-lg space-y-3 text-sm text-gray-200">
                            <div className="flex justify-between items-center border-b border-gray-600 pb-2">
                                <span className="text-gray-400">Тип</span>
                                <div className="flex items-center gap-2 font-medium">
                                    <OrgIcon size={16} className="text-[#DB935B]" />
                                    {org.org_type.toUpperCase()}
                                </div>
                            </div>
                            <div className="flex justify-between items-center border-b border-gray-600 pb-2">
                                <span className="text-gray-400">Баланс</span>
                                <span className="font-bold text-[#DB935B] text-lg">{org.balance.toLocaleString()} ₽</span>
                            </div>
                            <div className="flex justify-between items-center">
                                <span className="text-gray-400">Создана</span>
                                <span>{org.created_at ? new Date(org.created_at).toLocaleDateString() : '-'}</span>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Правая колонка - Форма / Данные */}
                <div className="lg:col-span-2">
                    <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-6 gap-4">
                        <div>
                            <h1 className="text-3xl font-bold text-white uppercase font-montserrat">Профиль организации</h1>
                            <p className="text-gray-300 text-sm mt-1">{orgTypeName}</p>
                        </div>

                        {!isEditing ? (
                            <Button
                                onClick={() => setIsEditing(true)}
                                className="bg-[#825e9c] text-white hover:bg-[#6a4c80] font-bold"
                            >
                                <Edit3 className="w-4 h-4 mr-2"/> Редактировать
                            </Button>
                        ) : (
                            <Button variant="outline" className="border-red-400 text-red-400 hover:bg-red-400/10" onClick={() => setIsEditing(false)}>
                                Отмена
                            </Button>
                        )}
                    </div>

                    <div className="bg-[#656662] p-1 rounded-lg shadow-xl border border-gray-500 relative">
                        {/* Форма редактирования/просмотра */}
                        <OrganisationForm
                            initialData={org}
                            onSubmit={handleUpdate}
                            isEditing={isEditing}
                        />

                        {/* Блокировка формы при просмотре (прозрачный слой) */}
                        {!isEditing && (
                            <div className="absolute inset-0 z-10 bg-transparent cursor-default" />
                        )}
                    </div>

                    <div className="mt-4">
                        <MessageComponent message={message}/>
                    </div>
                </div>
            </div>
        </div>
    );
}