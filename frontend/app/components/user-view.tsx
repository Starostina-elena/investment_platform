"use client";

import React, { useRef, useState, useEffect } from 'react';
import { Camera, MapPin, Mail, Edit3, Save, X } from 'lucide-react';
import { Button } from '@/app/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/app/components/ui/card';
import { Input } from "@/app/components/ui/input"; // Убедитесь, что компонент Input существует (shadcn)
import styles from '@/app/user-profile/page.module.css';
import {
    User,
    UpdateUserInfo,
    UploadUserAvatar,
    GetActiveInvestments,
    GetArchivedInvestments,
    Investment
} from "@/api/user";
import { useUserStore } from "@/context/user-store";
import { useImage } from "@/hooks/use-image";
import { BUCKETS } from "@/lib/config";
import Image from "next/image";
import avatarPlaceholder from "@/public/avatar.svg"; // Или любой другой плейсхолдер
import MessageComponent from "@/app/components/message";
import { Message } from "@/api/api";
import Spinner from "@/app/components/spinner";
import TopUpModal from "@/app/components/top-up-modal";
import BannedBanner from "@/app/components/banned-banner";
import AdminBanControl from "@/app/components/admin-ban-control";
import {Tabs, TabsContent, TabsList, TabsTrigger} from "@/app/components/ui/tabs";
import InvestmentsList from "@/app/components/investments-list";

interface UserViewProps {
    user: User;
    isOwner: boolean; // Флаг: это мой профиль?
}

export default function UserView({ user, isOwner }: UserViewProps) {
    const [isEditing, setIsEditing] = useState(false);

    // Стейт формы
    const [formData, setFormData] = useState<User>(user);

    // Стейт аватара
    const [avatarFile, setAvatarFile] = useState<File | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    // UI стейты
    const [message, setMessage] = useState<Message | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [activeInvestments, setActiveInvestments] = useState<Investment[]>([]);
    const [archivedInvestments, setArchivedInvestments] = useState<Investment[]>([]);
    const [loadingInvestments, setLoadingInvestments] = useState(false);

    // Загрузка инвестиций (только если это МОЙ профиль, т.к. инфа приватная)
    useEffect(() => {
        if (isOwner) {
            setLoadingInvestments(true);
            Promise.all([GetActiveInvestments(), GetArchivedInvestments()])
                .then(([active, archived]) => {
                    setActiveInvestments(active);
                    setArchivedInvestments(archived);
                })
                .finally(() => setLoadingInvestments(false));
        }
    }, [isOwner]);

    // Глобальный стор (чтобы обновить данные в хедере после сохранения)
    const { Login, token } = useUserStore();

    // Хук для картинки
    const avatarSrc = useImage(
        avatarFile,
        formData.avatar_path,
        BUCKETS.AVATARS,
        avatarPlaceholder
    );

    // Сброс формы при отмене
    const handleCancel = () => {
        setFormData(user);
        setAvatarFile(null);
        setIsEditing(false);
        setMessage(null);
    };

    const handleSave = async () => {
        setIsLoading(true);
        setMessage(null);

        // 1. Обновляем текстовые данные
        const updatedUser = await UpdateUserInfo(formData, setMessage);

        if (updatedUser) {
            // 2. Если выбран новый аватар, грузим его
            if (avatarFile) {
                const newAvatarPath = await UploadUserAvatar(avatarFile, setMessage);
                if (newAvatarPath) {
                    updatedUser.avatar_path = newAvatarPath;
                }
            }

            // 3. Обновляем глобальный стор
            if (token) {
                Login(updatedUser, token);
            }

            // 4. Обновляем локальный стейт
            setFormData(updatedUser);
            setIsEditing(false);
            setAvatarFile(null);
        }
        setIsLoading(false);
    };

    // Обработчик выбора файла
    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            if (file.size > 5 * 1024 * 1024) {
                setMessage({ isError: true, message: "Файл слишком большой (макс 5МБ)" });
                return;
            }
            setAvatarFile(file);
        }
    };

    return (
        <div className={styles.container}>
            {/* Баннер виден всем, если юзер забанен */}
            {formData.is_banned && (
                <div className="p-6">
                    <BannedBanner type="Пользователь" />
                </div>
            )}

            {/* Header Section */}
            <div className={styles.header}>
                <div className={styles.coverImage}>
                    <div className={styles.coverPattern}></div>
                </div>

                {/* Profile Section */}
                <div className={styles.profileSection}>
                    <div className={styles.profileGrid}>
                        {/* Left Profile Info */}
                        <div className={styles.leftColumn}>
                            <Card className={styles.profileCard}>
                                <CardContent className="p-6">
                                    {/* Profile Photo */}
                                    <div className={styles.profilePhoto}>
                                        <div className={styles.photoContainer}>
                                            <div className={styles.photoPlaceholder} style={{overflow: 'hidden'}}>
                                                <Image
                                                    src={avatarSrc}
                                                    alt={formData.nickname}
                                                    fill
                                                    style={{objectFit: 'cover'}}
                                                />
                                            </div>

                                            {/* Кнопка загрузки фото (только при редактировании) */}
                                            {isEditing && (
                                                <Button
                                                    variant="secondary"
                                                    size="sm"
                                                    className={styles.uploadAvatarBtn}
                                                    onClick={() => fileInputRef.current?.click()}
                                                >
                                                    <Camera className="w-4 h-4" />
                                                </Button>
                                            )}
                                            <input
                                                type="file"
                                                ref={fileInputRef}
                                                hidden
                                                accept="image/jpeg,image/png,image/webp"
                                                onChange={handleFileChange}
                                            />
                                        </div>
                                    </div>

                                    {/* Name */}
                                    <div className={styles.nameSection}>
                                        {isEditing ? (
                                            <div className="flex flex-col gap-2">
                                                <Input
                                                    value={formData.nickname}
                                                    onChange={e => setFormData({...formData, nickname: e.target.value})}
                                                    placeholder="Никнейм"
                                                />
                                            </div>
                                        ) : (
                                            <h1 className={styles.userName}>{formData.nickname}</h1>
                                        )}
                                        <p style={{color: '#aaa', fontSize: '0.9rem'}}>
                                            {formData.is_admin ? "Администратор" : "Пользователь"}
                                        </p>
                                    </div>

                                    {/* Contact Info */}
                                    <div className={styles.contactInfo}>
                                        <div className={styles.contactItem}>
                                            <span className={styles.contactLabel}>ID:</span>
                                            <span>{formData.id}</span>
                                        </div>
                                        <div className={styles.contactItem}>
                                            <span className={styles.contactLabel}>Баланс:</span>
                                            <span className="font-bold text-[#DB935B] text-lg">{formData.balance} ₽</span>
                                            {/* Показываем кнопку пополнения только владельцу */}
                                            {isOwner && <TopUpModal />}
                                        </div>

                                        <div className={styles.contactItem}>
                                            <Mail className="w-4 h-4 mr-2" />
                                            {isEditing ? (
                                                <Input
                                                    value={formData.email}
                                                    onChange={e => setFormData({...formData, email: e.target.value})}
                                                />
                                            ) : (
                                                <span>{formData.email}</span>
                                            )}
                                        </div>
                                        <div className={styles.contactItem}>
                                            <MapPin className="w-4 h-4 mr-2" />
                                            <span>Россия</span>
                                        </div>
                                    </div>

                                    {/* Edit / Save Buttons */}
                                    {isOwner && (
                                        <div className="mt-6 flex gap-2">
                                            {!isEditing ? (
                                                <Button className="w-full" onClick={() => setIsEditing(true)}>
                                                    <Edit3 className="w-4 h-4 mr-2" />
                                                    Редактировать
                                                </Button>
                                            ) : (
                                                <>
                                                    <Button variant="outline" onClick={handleCancel} disabled={isLoading}>
                                                        <X className="w-4 h-4" />
                                                    </Button>
                                                    <Button className="w-full" onClick={handleSave} disabled={isLoading}>
                                                        {isLoading ? <Spinner size={20} /> : <Save className="w-4 h-4 mr-2" />}
                                                        Сохранить
                                                    </Button>
                                                </>
                                            )}
                                        </div>
                                    )}

                                    {isOwner ? (
                                        <Card>
                                            <CardHeader>
                                                <CardTitle className={styles.descHeader}>Мои инвестиции</CardTitle>
                                            </CardHeader>
                                            <CardContent>
                                                {loadingInvestments ? (
                                                    <div className="flex justify-center p-4"><Spinner size={40} /></div>
                                                ) : (
                                                    <Tabs defaultValue="active" className="w-full">
                                                        <TabsList className="grid w-full grid-cols-2 bg-[#2d1b4e]">
                                                            <TabsTrigger value="active">Активные</TabsTrigger>
                                                            <TabsTrigger value="archived">Архив</TabsTrigger>
                                                        </TabsList>
                                                        <TabsContent value="active" className="mt-4">
                                                            <InvestmentsList investments={activeInvestments} />
                                                        </TabsContent>
                                                        <TabsContent value="archived" className="mt-4">
                                                            <InvestmentsList investments={archivedInvestments} />
                                                        </TabsContent>
                                                    </Tabs>
                                                )}
                                            </CardContent>
                                        </Card>
                                    ) : (
                                        // Если чужой профиль - показываем публичную инфу или заглушку
                                        <Card>
                                            <CardHeader>
                                                <CardTitle className={styles.descHeader}>Активность</CardTitle>
                                            </CardHeader>
                                            <CardContent>
                                                <p className={styles.emptyState}>Информация об инвестициях пользователя скрыта.</p>
                                            </CardContent>
                                        </Card>
                                    )}
                                    <AdminBanControl
                                        entityType="user"
                                        entityId={formData.id}
                                        isBanned={formData.is_banned}
                                        onUpdate={(banned) => setFormData({...formData, is_banned: banned})}
                                    />

                                    <MessageComponent message={message} style={{marginTop: '1rem'}}/>
                                </CardContent>
                            </Card>
                        </div>

                        {/* Main Content (Right Side) */}
                        <div className={styles.mainContent}>
                            <Card>
                                <CardHeader>
                                    <CardTitle className={styles.descHeader}>Личные данные</CardTitle>
                                </CardHeader>
                                <CardContent className="grid gap-4">
                                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                        <div>
                                            <label className="text-sm text-gray-500">Имя</label>
                                            {isEditing ? (
                                                <Input value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} />
                                            ) : <p className={styles.descText}>{formData.name}</p>}
                                        </div>
                                        <div>
                                            <label className="text-sm text-gray-500">Фамилия</label>
                                            {isEditing ? (
                                                <Input value={formData.surname} onChange={e => setFormData({...formData, surname: e.target.value})} />
                                            ) : <p className={styles.descText}>{formData.surname}</p>}
                                        </div>
                                        <div>
                                            <label className="text-sm text-gray-500">Отчество</label>
                                            {isEditing ? (
                                                <Input value={formData.patronymic || ''} onChange={e => setFormData({...formData, patronymic: e.target.value})} />
                                            ) : <p className={styles.descText}>{formData.patronymic || '-'}</p>}
                                        </div>
                                    </div>

                                    <div className="mt-4">
                                        <label className="text-sm text-gray-500">Дата регистрации</label>
                                        <p className={styles.descText}>
                                            {new Date(formData.created_at).toLocaleDateString('ru-RU')}
                                        </p>
                                    </div>
                                </CardContent>
                            </Card>

                            <Card>
                                <CardHeader>
                                    <CardTitle className={styles.descHeader}>Активность</CardTitle>
                                </CardHeader>
                                <CardContent>
                                    <p className={styles.emptyState}>История инвестиций и проектов будет здесь.</p>
                                </CardContent>
                            </Card>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}