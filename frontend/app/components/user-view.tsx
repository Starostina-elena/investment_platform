"use client";

import React, { useRef, useState, useEffect } from 'react';
import { Camera, MapPin, Mail, Edit3, Save, X, Hash, Wallet, LogOut, Building2 } from 'lucide-react'; // Добавил иконки
import { Button } from '@/app/components/ui/button';
import { Input } from "@/app/components/ui/input";
import styles from '@/app/user-profile/page.module.css';
import { User, UpdateUserInfo, UploadUserAvatar, GetActiveInvestments, GetArchivedInvestments, Investment } from "@/api/user";
import { useUserStore } from "@/context/user-store";
import { useImage } from "@/hooks/use-image";
import { BUCKETS } from "@/lib/config";
import Image from "next/image";
import avatarPlaceholder from "@/public/avatar.svg";
import MessageComponent from "@/app/components/message";
import { Message } from "@/api/api";
import Spinner from "@/app/components/spinner";
import TopUpModal from "@/app/components/top-up-modal";
import InvestmentsList from "@/app/components/investments-list";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/app/components/ui/tabs";
import AdminBanControl from "@/app/components/admin-ban-control";
import BannedBanner from "@/app/components/banned-banner";
import Link from "next/link";
import { useRouter } from "next/navigation";

interface UserViewProps {
    user: User;
    isOwner: boolean;
}

export default function UserView({ user, isOwner }: UserViewProps) {
    const [isEditing, setIsEditing] = useState(false);
    const [formData, setFormData] = useState<User>(user);
    const [avatarFile, setAvatarFile] = useState<File | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);
    const [message, setMessage] = useState<Message | null>(null);
    const [isLoading, setIsLoading] = useState(false);

    // Получаем Logout
    const { Login, Logout, token } = useUserStore();
    const router = useRouter();

    // Стейт для инвестиций
    const [activeInvestments, setActiveInvestments] = useState<Investment[]>([]);
    const [archivedInvestments, setArchivedInvestments] = useState<Investment[]>([]);
    const [loadingInvestments, setLoadingInvestments] = useState(false);

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

    const avatarSrc = useImage(
        avatarFile,
        formData.avatar_path,
        BUCKETS.AVATARS,
        avatarPlaceholder
    );

    const handleCancel = () => {
        setFormData(user);
        setAvatarFile(null);
        setIsEditing(false);
        setMessage(null);
    };

    const handleSave = async () => {
        setIsLoading(true);
        const responseData = await UpdateUserInfo(formData, setMessage);

        if (responseData) {
            let newAvatarPath = user.avatar_path; // По умолчанию старая аватарка

            // Если загружали новую аву
            if (avatarFile) {
                const uploadedPath = await UploadUserAvatar(avatarFile, setMessage);
                if (uploadedPath) newAvatarPath = uploadedPath;
            }

            // МЕРДЖИМ ДАННЫЕ:
            // Берем исходного юзера (с правильным балансом и датой)
            // И перезаписываем только те поля, которые мы редактировали
            const mergedUser: User = {
                ...user, // <- Основа (ID, Balance, CreatedAt, IsAdmin, IsBanned)

                // Обновляемые поля берем из ответа сервера (или из формы)
                name: responseData.name,
                surname: responseData.surname,
                patronymic: responseData.patronymic,
                nickname: responseData.nickname,
                email: responseData.email,

                // Аватар обновляем отдельно
                avatar_path: newAvatarPath
            };

            // Обновляем глобальный стор
            if (token) Login(mergedUser, token);

            // Обновляем локальную форму
            setFormData(mergedUser);
            setIsEditing(false);
            setAvatarFile(null);
        }
        setIsLoading(false);
    };

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) setAvatarFile(file);
    };

    const handleLogout = () => {
        Logout();
        router.push("/login");
    };

    return (
        <div className={styles.container}>
            {/* Баннер бана */}
            {formData.is_banned && (
                <div className="max-w-6xl mx-auto pt-6 px-8">
                    <BannedBanner type="Пользователь" />
                </div>
            )}

            {/* 1. Обложка */}
            <div className={styles.header}>
                <div className={styles.coverImage}>
                    <div className={styles.coverPattern}></div>
                </div>

                {/* 2. Секция профиля */}
                <div className={styles.profileSection}>
                    <div className={styles.profileGrid}>

                        {/* Левая колонка: Карточка */}
                        <div className={styles.leftColumn}>
                            <div className={styles.profileCard}>
                                {/* Фото */}
                                <div className={styles.profilePhoto}>
                                    <div className={styles.photoContainer}>
                                        <div className={styles.photoPlaceholder} style={{position: 'relative', background: '#fff'}}>
                                            {/* Добавил unoptimized, чтобы избежать проблем с доменами */}
                                            <Image
                                                src={avatarSrc}
                                                alt={formData.nickname}
                                                fill
                                                unoptimized
                                                style={{objectFit: 'cover'}}
                                            />
                                        </div>
                                        {isEditing && (
                                            <button
                                                className={styles.uploadAvatarBtn}
                                                onClick={() => fileInputRef.current?.click()}
                                                type="button"
                                            >
                                                <Camera size={18} />
                                            </button>
                                        )}
                                        <input type="file" ref={fileInputRef} hidden onChange={handleFileChange} accept="image/*" />
                                    </div>
                                </div>

                                {/* Имя */}
                                <div className={styles.nameSection}>
                                    {isEditing ? (
                                        <Input
                                            className="text-center font-bold text-lg mb-2 text-black bg-white border-gray-300"
                                            value={formData.nickname}
                                            onChange={e => setFormData({...formData, nickname: e.target.value})}
                                            placeholder="Никнейм"
                                        />
                                    ) : (
                                        <h1 className={styles.userName}>{formData.nickname}</h1>
                                    )}
                                    <div className={styles.userStatus}>
                                        {formData.is_admin ? "Администратор" : "Пользователь"}
                                    </div>
                                </div>

                                {/* Инфо */}
                                <div className={styles.contactInfo}>
                                    <div className={styles.contactItem}>
                                        <span className={styles.contactLabel}><Hash size={14} className="inline mr-1"/>ID</span>
                                        <span className="font-mono font-bold">{formData.id}</span>
                                    </div>

                                    <div className={styles.contactItem}>
                                        <span className={styles.contactLabel}><Wallet size={14} className="inline mr-1"/>Баланс</span>
                                        <div className="flex items-center gap-2">
                                            <span className="font-bold text-[#DB935B] text-lg">{(formData.balance ?? 0).toLocaleString()} ₽</span>
                                            {isOwner && <TopUpModal />}
                                        </div>
                                    </div>

                                    <div className={styles.contactItem}>
                                        <span className={styles.contactLabel}><Mail size={14} className="inline mr-1"/>Email</span>
                                        {isEditing ? (
                                            <Input
                                                className="h-8 w-40 text-right text-xs text-black bg-white border-gray-300"
                                                value={formData.email}
                                                onChange={e => setFormData({...formData, email: e.target.value})}
                                            />
                                        ) : (
                                            <span className="truncate max-w-[150px] font-medium" title={formData.email}>{formData.email}</span>
                                        )}
                                    </div>

                                    <div className={styles.contactItem}>
                                        <span className={styles.contactLabel}><MapPin size={14} className="inline mr-1"/>Страна</span>
                                        <span className="font-medium">Россия</span>
                                    </div>
                                </div>

                                {/* Кнопки Владельца */}
                                {isOwner && (
                                    <div className="px-6 mt-6 pb-4 flex flex-col gap-3">
                                        {!isEditing ? (
                                            <>
                                                {/* Редактировать */}
                                                <Button
                                                    className="w-full bg-[#656662] hover:bg-[#505050] text-white font-bold uppercase tracking-wider"
                                                    onClick={() => setIsEditing(true)}
                                                >
                                                    <Edit3 className="w-4 h-4 mr-2" /> Редактировать
                                                </Button>

                                                {/* Создать организацию */}
                                                <Link href="/organisation/create" className="w-full">
                                                    <Button variant="outline" className="w-full border-[#825e9c] text-[#825e9c] hover:bg-[#825e9c] hover:text-white font-bold uppercase tracking-wider">
                                                        <Building2 className="w-4 h-4 mr-2" /> Создать орг-цию
                                                    </Button>
                                                </Link>

                                                {/* Мои организации */}
                                                <Link href="/organisation/my" className="w-full">
                                                    <Button variant="outline" className="w-full border-gray-400 text-gray-600 hover:bg-gray-100 font-bold uppercase tracking-wider">
                                                        Мои организации
                                                    </Button>
                                                </Link>

                                                {/* Выйти */}
                                                <Button
                                                    variant="ghost"
                                                    className="w-full text-red-500 hover:bg-red-50 hover:text-red-700 mt-2"
                                                    onClick={handleLogout}
                                                >
                                                    <LogOut className="w-4 h-4 mr-2" /> Выйти
                                                </Button>
                                            </>
                                        ) : (
                                            <div className="flex gap-2">
                                                <Button variant="outline" onClick={handleCancel} className="flex-1 border-gray-400 hover:bg-gray-100 text-black">
                                                    <X className="w-4 h-4" />
                                                </Button>
                                                <Button
                                                    onClick={handleSave}
                                                    disabled={isLoading}
                                                    className="flex-1 bg-[#825e9c] hover:bg-[#6a4c80] text-white font-bold"
                                                >
                                                    {isLoading ? <Spinner size={16} /> : <Save className="w-4 h-4 mr-2" />} Сохр.
                                                </Button>
                                            </div>
                                        )}
                                        <MessageComponent message={message} style={{marginTop: '0.5rem', fontSize: '0.8rem'}}/>

                                        {/* Админ панель */}
                                        <AdminBanControl
                                            entityType="user"
                                            entityId={formData.id}
                                            isBanned={formData.is_banned}
                                            onUpdate={(banned) => setFormData({...formData, is_banned: banned})}
                                        />
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* Правая колонка: Данные */}
                        <div className={styles.mainContent}>

                            {/* Личные данные */}
                            <div>
                                <h2 className={styles.descHeader}>Личные данные</h2>
                                <div className={styles.infoGrid}>
                                    <div className={styles.infoItem}>
                                        <label>Имя</label>
                                        {isEditing ? (
                                            <Input className="bg-gray-100 text-black border-gray-300" value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} />
                                        ) : <p>{formData.name}</p>}
                                    </div>
                                    <div className={styles.infoItem}>
                                        <label>Фамилия</label>
                                        {isEditing ? (
                                            <Input className="bg-gray-100 text-black border-gray-300" value={formData.surname} onChange={e => setFormData({...formData, surname: e.target.value})} />
                                        ) : <p>{formData.surname}</p>}
                                    </div>
                                    <div className={styles.infoItem}>
                                        <label>Отчество</label>
                                        {isEditing ? (
                                            <Input className="bg-gray-100 text-black border-gray-300" value={formData.patronymic || ''} onChange={e => setFormData({...formData, patronymic: e.target.value})} />
                                        ) : <p>{formData.patronymic || '—'}</p>}
                                    </div>
                                    <div className={styles.infoItem}>
                                        <label>Дата регистрации</label>
                                        <p>{new Date(formData.created_at).toLocaleDateString('ru-RU')}</p>
                                    </div>
                                </div>
                            </div>

                            {/* Активность / Инвестиции */}
                            <div>
                                <h2 className={styles.descHeader}>Мои инвестиции</h2>
                                {isOwner ? (
                                    <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
                                        {loadingInvestments ? (
                                            <div className="flex justify-center p-8"><Spinner size={40} /></div>
                                        ) : (
                                            <Tabs defaultValue="active" className="w-full">
                                                <TabsList className="grid w-full grid-cols-2 bg-gray-200 mb-6">
                                                    <TabsTrigger value="active" className="data-[state=active]:bg-[#825e9c] data-[state=active]:text-white font-bold">Активные</TabsTrigger>
                                                    <TabsTrigger value="archived" className="data-[state=active]:bg-[#656662] data-[state=active]:text-white font-bold">Архив</TabsTrigger>
                                                </TabsList>
                                                <TabsContent value="active">
                                                    <InvestmentsList investments={activeInvestments} />
                                                </TabsContent>
                                                <TabsContent value="archived">
                                                    <InvestmentsList investments={archivedInvestments} />
                                                </TabsContent>
                                            </Tabs>
                                        )}
                                    </div>
                                ) : (
                                    <div className="bg-white p-6 rounded-lg shadow-sm text-center">
                                        <p className="text-gray-500 italic">Информация об инвестициях пользователя скрыта.</p>
                                    </div>
                                )}
                            </div>

                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}