'use client'

import React, {useEffect, useRef, useState} from 'react';
import {useParams} from "next/navigation";
import {GetFullOrganisation, Organisation, UpdateOrganisation, UploadOrgAvatar} from "@/api/organisation";
import OrganisationForm from "@/app/components/organisation-form";
import Spinner from "@/app/components/spinner";
import {useImage} from "@/hooks/use-image";
import {BUCKETS} from "@/lib/config";
import Image from "next/image";
import {Button} from "@/app/components/ui/button";
import {Camera, Edit3, Building2, UserCircle, Briefcase} from "lucide-react";
import MessageComponent from "@/app/components/message";
import {Message} from "@/api/api";
import {Badge} from "@/app/components/ui/badge";
import orgPlaceholder from "@/public/image_bg.png";

export default function OrganisationPage() {
    const params = useParams();
    const [org, setOrg] = useState<Organisation | null | undefined>(undefined);
    const [isEditing, setIsEditing] = useState(false);
    const [message, setMessage] = useState<Message | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

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
            setOrg({...org, ...updated});
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

    const OrgIcon = org.org_type === 'ip' ? Briefcase : org.org_type === 'jur' ? Building2 : UserCircle;
    const orgTypeName = org.org_type === 'ip' ? 'Индивидуальный предприниматель' : org.org_type === 'jur' ? 'Юридическое лицо' : 'Физическое лицо';

    return (
        <div style={{minHeight: '100vh', backgroundColor: '#989694', padding: '3rem 2rem'}}>
            <div style={{maxWidth: '1200px', margin: '0 auto', display: 'grid', gridTemplateColumns: '320px 1fr', gap: '2rem'}}>

                {/* Левая колонка - Карточка */}
                <div style={{width: '100%'}}>
                    <div style={{
                        backgroundColor: '#656662',
                        padding: '2rem',
                        borderRadius: '8px',
                        boxShadow: '0 4px 20px rgba(0,0,0,0.2)',
                        border: '1px solid #4a4a4a',
                        textAlign: 'center',
                        position: 'sticky',
                        top: '100px',
                        overflow: 'hidden' // ВАЖНО: чтобы ничего не вылезало
                    }}>
                        {/* Аватар: Жесткий контейнер */}
                        <div style={{
                            width: '160px',
                            height: '160px',
                            margin: '0 auto 1.5rem auto',
                            borderRadius: '50%',
                            overflow: 'hidden',
                            border: '4px solid #825e9c',
                            backgroundColor: 'white',
                            cursor: 'default',
                            pointerEvents: 'none'
                        }}>
                            <img
                                src={typeof avatarSrc === 'string' ? avatarSrc : avatarSrc.src}
                                alt={org.name}
                                style={{width: '100%', height: '100%', objectFit: 'cover', pointerEvents: 'none', userSelect: 'none'}}
                                draggable={false}
                            />
                        </div>

                        {!isEditing && (
                            <>
                                <input type="file" ref={fileInputRef} hidden onChange={handleAvatarUpload} accept="image/*"/>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    style={{marginBottom: '1.5rem', color: 'white', borderColor: '#aaa', backgroundColor: 'transparent'}}
                                    onClick={() => fileInputRef.current?.click()}
                                >
                                    <Camera size={16} style={{marginRight: '8px'}}/> Сменить логотип
                                </Button>
                            </>
                        )}

                        <h2 style={{fontSize: '1.5rem', fontWeight: '800', color: 'white', marginBottom: '0.5rem', fontFamily: 'var(--font-montserrat)', wordWrap: 'break-word'}}>
                            {org.name}
                        </h2>

                        <div style={{display: 'flex', justifyContent: 'center', gap: '10px', marginBottom: '1.5rem'}}>
                            <Badge className={`${org.registration_completed ? 'bg-green-600' : 'bg-yellow-600'} text-white border-none`}>
                                {org.registration_completed ? 'Активна' : 'Черновик'}
                            </Badge>
                            {org.is_banned && <Badge variant="destructive">Заблокирована</Badge>}
                        </div>

                        <div style={{textAlign: 'left', backgroundColor: '#555652', padding: '1rem', borderRadius: '8px', color: '#e0e0e0', fontSize: '0.9rem'}}>
                            <div style={{display: 'flex', justifyContent: 'space-between', paddingBottom: '0.5rem', borderBottom: '1px solid #666', marginBottom: '0.5rem'}}>
                                <span style={{color: '#aaa'}}>Тип</span>
                                <div style={{display: 'flex', alignItems: 'center', gap: '5px', fontWeight: 'bold'}}>
                                    <OrgIcon size={16} color="#DB935B" />
                                    <span style={{textTransform: 'uppercase'}}>{org.org_type}</span>
                                </div>
                            </div>
                            <div style={{display: 'flex', justifyContent: 'space-between', paddingBottom: '0.5rem', borderBottom: '1px solid #666', marginBottom: '0.5rem'}}>
                                <span style={{color: '#aaa'}}>Баланс</span>
                                <span style={{color: '#DB935B', fontWeight: 'bold', fontSize: '1.1rem'}}>
                                    {(org.balance ?? 0).toLocaleString()} ₽
                                </span>
                            </div>
                            <div style={{display: 'flex', justifyContent: 'space-between'}}>
                                <span style={{color: '#aaa'}}>Создана</span>
                                <span>{org.created_at ? new Date(org.created_at).toLocaleDateString() : '-'}</span>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Правая колонка - Данные */}
                <div style={{minWidth: 0}}> {/* minWidth: 0 важен для Grid, чтобы контент не распирал колонку */}
                    <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem', flexWrap: 'wrap', gap: '1rem'}}>
                        <div>
                            <h1 style={{fontSize: '2rem', fontWeight: '800', color: 'white', textTransform: 'uppercase', fontFamily: 'var(--font-montserrat)'}}>
                                Профиль организации
                            </h1>
                            <p style={{color: '#ccc'}}>{orgTypeName}</p>
                        </div>

                        {!isEditing ? (
                            <Button
                                onClick={() => setIsEditing(true)}
                                style={{backgroundColor: '#825e9c', color: 'white', fontWeight: 'bold'}}
                            >
                                <Edit3 size={16} style={{marginRight: '8px'}}/> Редактировать
                            </Button>
                        ) : (
                            <Button variant="outline" onClick={() => setIsEditing(false)} style={{color: '#ff9999', borderColor: '#ff9999'}}>
                                Отмена
                            </Button>
                        )}
                    </div>

                    <div style={{backgroundColor: '#656662', padding: '2px', borderRadius: '8px', boxShadow: '0 4px 10px rgba(0,0,0,0.1)', position: 'relative'}}>
                        <OrganisationForm
                            initialData={org}
                            onSubmit={handleUpdate}
                            isEditing={isEditing}
                        />

                        {!isEditing && (
                            <div style={{position: 'absolute', inset: 0, zIndex: 10, backgroundColor: 'transparent'}} />
                        )}
                    </div>

                    <div style={{marginTop: '1rem'}}>
                        <MessageComponent message={message}/>
                    </div>
                </div>
            </div>

            <style jsx>{`
                @media (max-width: 1024px) {
                    div[style*="gridTemplateColumns"] {
                        grid-template-columns: 1fr !important;
                    }
                }
            `}</style>
        </div>
    );
}