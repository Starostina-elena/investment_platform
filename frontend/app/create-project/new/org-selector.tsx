'use client'
import React from 'react';
import {Organisation} from "@/api/organisation";
import {Project} from "@/api/project";
import {Card, CardContent, CardHeader, CardTitle} from "@/app/components/ui/card";
import {Button} from "@/app/components/ui/button";
import Link from "next/link";
import styles from "./generic-info.module.css"; // Используем существующие стили

interface OrgSelectorProps {
    project: Project;
    setProject: (p: Project) => void;
    userOrgs: Organisation[];
    children?: React.ReactNode; // Для кнопок навигации
}

export default function OrgSelector({project, setProject, userOrgs, children}: OrgSelectorProps) {

    if (userOrgs.length === 0) {
        return (
            <div className={styles.main_content} style={{textAlign: 'center'}}>
                <h2 className={styles.page_title}>Сначала создайте организацию</h2>
                <p className={styles.welcome_text}>
                    Чтобы запустить проект, вам необходимо юридическое лицо (или статус ИП/Самозанятого),
                    к которому будут привязаны банковские реквизиты для сбора средств.
                </p>
                <Link href="/organisation/create">
                    <Button className="bg-light-blue text-black font-bold text-lg px-8 py-6">
                        Зарегистрировать организацию
                    </Button>
                </Link>
            </div>
        );
    }

    return (
        <div className={styles.main_content}>
            <h1 className={styles.page_title}>Выбор автора проекта</h1>
            <p className={styles.welcome_text}>
                Выберите организацию, от лица которой будет опубликован этот проект.
                Средства будут поступать на счет этой организации.
            </p>

            <div className="grid gap-4">
                {userOrgs.map(org => (
                    <Card
                        key={org.id}
                        className={`cursor-pointer transition-all border-2 ${project.creator_id === org.id ? 'border-light-blue bg-[#f0f9ff]' : 'border-transparent hover:border-gray-300'}`}
                        onClick={() => setProject({...project, creator_id: org.id})}
                    >
                        <CardHeader className="pb-2">
                            <CardTitle className="flex justify-between items-center text-black">
                                {org.name}
                                {project.creator_id === org.id && <span className="text-light-blue text-sm">Выбрано</span>}
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <p className="text-sm text-gray-500">
                                Тип: {org.org_type.toUpperCase()} | Баланс: {org.balance} ₽
                            </p>
                            {!org.registration_completed && (
                                <p className="text-red-500 text-xs mt-1">
                                    ⚠ Реквизиты не заполнены. Вы не сможете вывести средства.
                                    <Link href={`/organisation/${org.id}`} className="underline ml-1">Заполнить</Link>
                                </p>
                            )}
                        </CardContent>
                    </Card>
                ))}
            </div>

            <div className="mt-4">
                <Link href="/organisation/create" className="text-sm text-gray-400 underline">
                    + Создать новую организацию
                </Link>
            </div>

            {/* Рендерим кнопки навигации (Вперед/Назад), только если выбрана организация */}
            {project.creator_id > 0 && children}
        </div>
    );
}