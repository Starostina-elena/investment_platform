'use client'
import React from 'react';
import {Organisation} from "@/api/organisation";
import {Project} from "@/api/project";
import {Card, CardContent, CardHeader, CardTitle} from "@/app/components/ui/card";
import {Button} from "@/app/components/ui/button";
import Link from "next/link";
import { CheckCircle, AlertTriangle } from "lucide-react";

interface OrgSelectorProps {
    project: Project;
    setProject: (p: Project) => void;
    userOrgs: Organisation[];
    children?: React.ReactNode;
}

export default function OrgSelector({project, setProject, userOrgs, children}: OrgSelectorProps) {

    if (userOrgs.length === 0) {
        return (
            <div className="bg-[#656662] p-10 rounded-lg text-center shadow-xl text-white">
                <h2 className="text-3xl font-bold uppercase mb-4">Сначала создайте организацию</h2>
                <p className="text-gray-300 mb-8 max-w-2xl mx-auto">
                    Чтобы запустить проект, вам необходимо юридическое лицо (или статус ИП/Самозанятого),
                    к которому будут привязаны банковские реквизиты для сбора средств.
                </p>
                <Link href="/organisation/create">
                    <Button className="bg-[#825e9c] text-white hover:bg-[#6a4c80] font-bold text-lg px-8 py-6 uppercase tracking-wider shadow-lg">
                        Зарегистрировать организацию
                    </Button>
                </Link>
            </div>
        );
    }

    return (
        <div className="bg-[#656662] p-8 rounded-lg shadow-xl text-white">
            <h1 className="text-3xl font-bold uppercase mb-4 border-b border-gray-500 pb-4">Выбор автора проекта</h1>
            <p className="text-gray-300 mb-8">
                Выберите организацию, от лица которой будет опубликован этот проект.
                Средства будут поступать на счет этой организации.
            </p>

            <div className="grid gap-4 mb-8">
                {userOrgs.map(org => {
                    const isSelected = project.creator_id === org.id;
                    return (
                        <Card
                            key={org.id}
                            className={`cursor-pointer transition-all border-2 
                                ${isSelected
                                ? 'border-[#825e9c] bg-[#555652]'
                                : 'border-transparent bg-[#505050] hover:bg-[#5a5a5a]'
                            }`}
                            onClick={() => setProject({...project, creator_id: org.id})}
                        >
                            <CardHeader className="pb-2">
                                <div className="flex justify-between items-center w-full">
                                    <CardTitle className="text-white text-xl font-bold uppercase">
                                        {org.name + " "}
                                    </CardTitle>

                                    {/* Бейдж "Выбрано" справа */}
                                    {isSelected && (
                                        <div className="flex items-center gap-2 bg-[#825e9c] text-white px-3 py-1 rounded-full text-xs font-bold uppercase">
                                            <CheckCircle size={14} />
                                            Выбрано
                                        </div>
                                    )}
                                </div>
                            </CardHeader>
                            <CardContent>
                                <div className="flex items-center gap-4 text-sm text-gray-300 mb-2">
                                    <span className="bg-gray-700 px-2 py-1 rounded uppercase font-bold text-xs">
                                        {org.org_type}
                                    </span>
                                    <span>
                                        Баланс: <span className="font-bold text-[#DB935B]">{(org.balance ?? 0)} ₽</span>
                                    </span>
                                </div>

                                {!org.registration_completed && (
                                    <div className="flex items-center gap-2 text-red-400 text-xs bg-red-500/10 p-2 rounded border border-red-500/20 mt-2">
                                        <AlertTriangle size={14} />
                                        <span>Реквизиты не заполнены. Вы не сможете вывести средства.</span>
                                        <Link href={`/organisation/${org.id}`} className="underline font-bold hover:text-red-300 ml-auto">
                                            Заполнить
                                        </Link>
                                    </div>
                                )}
                            </CardContent>
                        </Card>
                    );
                })}
            </div>

            <div className="flex justify-between items-center mt-6 pt-6 border-t border-gray-500">
                <Link href="/organisation/create" className="text-sm text-[#DB935B] hover:text-white transition-colors underline">
                    + Создать новую организацию
                </Link>

                {/* Кнопки навигации */}
                <div className="flex gap-4">
                    {project.creator_id > 0 && children}
                </div>
            </div>
        </div>
    );
}