'use client'

import React, { useState } from 'react';
import { Organisation, OrgType } from "@/api/organisation";
import { Input } from "@/app/components/ui/input";
import { Button } from "@/app/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/app/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/app/components/ui/card";
import { Label } from "@/app/components/ui/label";

interface OrganisationFormProps {
    initialData?: Partial<Organisation>;
    onSubmit: (data: Partial<Organisation>) => void;
    isEditing?: boolean;
}

export default function OrganisationForm({ initialData, onSubmit, isEditing = false }: OrganisationFormProps) {
    const [formData, setFormData] = useState<Partial<Organisation>>({
        org_type: 'phys',
        phys_face: {} as any,
        jur_face: {} as any,
        ip_face: {} as any,
        ...initialData
    });

    const handleChange = (field: string, value: any) => {
        setFormData(prev => ({ ...prev, [field]: value }));
    };

    const handleFaceChange = (faceType: 'phys_face' | 'jur_face' | 'ip_face', field: string, value: any) => {
        setFormData(prev => ({
            ...prev,
            [faceType]: {
                ...prev[faceType],
                [field]: value
            }
        }));
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        const cleanData = { ...formData };
        if (cleanData.org_type === 'phys') { delete cleanData.jur_face; delete cleanData.ip_face; }
        if (cleanData.org_type === 'jur') { delete cleanData.phys_face; delete cleanData.ip_face; }
        if (cleanData.org_type === 'ip') { delete cleanData.phys_face; delete cleanData.jur_face; }

        onSubmit(cleanData);
    };

    // Общие стили для инпутов, чтобы они были видны на темном фоне
    const inputClass = "bg-white text-black border-gray-400 placeholder:text-gray-500";
    const labelClass = "text-gray-300 mb-1 block";

    return (
        <form onSubmit={handleSubmit} className="space-y-8 max-w-4xl mx-auto">
            {/* Секция 1: Основные данные */}
            <Card className="bg-[#656662] border-gray-500 shadow-xl">
                <CardHeader>
                    <CardTitle className="text-white uppercase font-bold text-xl border-b border-gray-500 pb-2">Основные данные</CardTitle>
                </CardHeader>
                <CardContent className="grid gap-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="space-y-2">
                            <Label className={labelClass}>Название (для отображения)</Label>
                            <Input
                                className={inputClass}
                                required
                                value={formData.name || ''}
                                onChange={e => handleChange('name', e.target.value)}
                                placeholder='Например: "Мой Стартап"'
                            />
                        </div>
                        <div className="space-y-2">
                            <Label className={labelClass}>Email организации</Label>
                            <Input
                                className={inputClass}
                                required type="email"
                                value={formData.email || ''}
                                onChange={e => handleChange('email', e.target.value)}
                                placeholder="corp@example.com"
                            />
                        </div>
                    </div>

                    {!isEditing && (
                        <div className="space-y-2">
                            <Label className={labelClass}>Тип юридического лица</Label>
                            <Select
                                value={formData.org_type}
                                onValueChange={(val: OrgType) => handleChange('org_type', val)}
                            >
                                <SelectTrigger className={`${inputClass} w-full`}>
                                    <SelectValue placeholder="Выберите тип" />
                                </SelectTrigger>
                                <SelectContent className="bg-white text-black">
                                    <SelectItem value="phys">Физическое лицо</SelectItem>
                                    <SelectItem value="jur">Юридическое лицо (ООО, АО)</SelectItem>
                                    <SelectItem value="ip">Индивидуальный предприниматель</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>
                    )}
                </CardContent>
            </Card>

            {/* Секция 2: Реквизиты */}
            <Card className="bg-[#656662] border-gray-500 shadow-xl">
                <CardHeader>
                    <CardTitle className="text-white uppercase font-bold text-xl border-b border-gray-500 pb-2">Реквизиты</CardTitle>
                </CardHeader>
                <CardContent className="grid gap-6 pt-4">
                    {/* Контейнер грида для полей */}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">

                        {formData.org_type === 'phys' && (
                            <>
                                <div className="space-y-2">
                                    <Label className={labelClass}>ФИО</Label>
                                    <Input className={inputClass} placeholder="Иванов Иван Иванович" required
                                           value={formData.phys_face?.fio || ''}
                                           onChange={e => handleFaceChange('phys_face', 'fio', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>ИНН</Label>
                                    <Input className={inputClass} placeholder="12 цифр" required
                                           value={formData.phys_face?.inn || ''}
                                           onChange={e => handleFaceChange('phys_face', 'inn', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>Серия паспорта</Label>
                                    <Input className={inputClass} placeholder="1234" type="number" required
                                           value={formData.phys_face?.passport_series || ''}
                                           onChange={e => handleFaceChange('phys_face', 'passport_series', +e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>Номер паспорта</Label>
                                    <Input className={inputClass} placeholder="567890" type="number" required
                                           value={formData.phys_face?.passport_number || ''}
                                           onChange={e => handleFaceChange('phys_face', 'passport_number', +e.target.value)} />
                                </div>
                                <div className="space-y-2 md:col-span-2">
                                    <Label className={labelClass}>Кем выдан</Label>
                                    <Input className={inputClass} placeholder="Отделением УФМС..." required
                                           value={formData.phys_face?.passport_givenby || ''}
                                           onChange={e => handleFaceChange('phys_face', 'passport_givenby', e.target.value)} />
                                </div>

                                <div className="md:col-span-2 border-t border-gray-500 my-2"></div>

                                <div className="space-y-2">
                                    <Label className={labelClass}>БИК Банка</Label>
                                    <Input className={inputClass} placeholder="044..." required
                                           value={formData.phys_face?.bic || ''}
                                           onChange={e => handleFaceChange('phys_face', 'bic', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>Расчетный счет</Label>
                                    <Input className={inputClass} placeholder="408..." required
                                           value={formData.phys_face?.checking_account || ''}
                                           onChange={e => handleFaceChange('phys_face', 'checking_account', e.target.value)} />
                                </div>
                                <div className="space-y-2 md:col-span-2">
                                    <Label className={labelClass}>Корреспондентский счет</Label>
                                    <Input className={inputClass} placeholder="301..." required
                                           value={formData.phys_face?.correspondent_account || ''}
                                           onChange={e => handleFaceChange('phys_face', 'correspondent_account', e.target.value)} />
                                </div>

                                <div className="md:col-span-2 border-t border-gray-500 my-2"></div>

                                <div className="space-y-2 md:col-span-2">
                                    <Label className={labelClass}>Адрес регистрации</Label>
                                    <Input className={inputClass} placeholder="г. Москва, ул..." required
                                           value={formData.phys_face?.registration_address || ''}
                                           onChange={e => handleFaceChange('phys_face', 'registration_address', e.target.value)} />
                                </div>
                                <div className="space-y-2 md:col-span-2">
                                    <Label className={labelClass}>Почтовый адрес</Label>
                                    <Input className={inputClass} placeholder="Индекс, г. Москва..." required
                                           value={formData.phys_face?.post_address || ''}
                                           onChange={e => handleFaceChange('phys_face', 'post_address', e.target.value)} />
                                </div>
                            </>
                        )}

                        {/* Аналогично оформляем JurFace и IPFace, добавляя Label и обертки div */}
                        {formData.org_type === 'jur' && (
                            <>
                                <div className="space-y-2 md:col-span-2">
                                    <Label className={labelClass}>Полное наименование</Label>
                                    <Input className={inputClass} required
                                           value={formData.jur_face?.full_organisation_name || ''}
                                           onChange={e => handleFaceChange('jur_face', 'full_organisation_name', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>Краткое наименование</Label>
                                    <Input className={inputClass} required
                                           value={formData.jur_face?.short_organisation_name || ''}
                                           onChange={e => handleFaceChange('jur_face', 'short_organisation_name', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>ИНН</Label>
                                    <Input className={inputClass} required
                                           value={formData.jur_face?.inn || ''}
                                           onChange={e => handleFaceChange('jur_face', 'inn', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>ОГРН</Label>
                                    <Input className={inputClass} required
                                           value={formData.jur_face?.ogrn || ''}
                                           onChange={e => handleFaceChange('jur_face', 'ogrn', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>КПП</Label>
                                    <Input className={inputClass} required
                                           value={formData.jur_face?.kpp || ''}
                                           onChange={e => handleFaceChange('jur_face', 'kpp', e.target.value)} />
                                </div>
                                {/* ... остальные поля юр лица аналогично ... */}
                                {/* Чтобы не раздувать ответ, принцип тот же: Label + Input внутри div.space-y-2 */}
                                <div className="space-y-2">
                                    <Label className={labelClass}>Должность руководителя</Label>
                                    <Input className={inputClass} required value={formData.jur_face?.position || ''} onChange={e => handleFaceChange('jur_face', 'position', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>Действует на основании</Label>
                                    <Input className={inputClass} required value={formData.jur_face?.acts_on_base || ''} onChange={e => handleFaceChange('jur_face', 'acts_on_base', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>БИК</Label>
                                    <Input className={inputClass} required value={formData.jur_face?.bic || ''} onChange={e => handleFaceChange('jur_face', 'bic', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>Расчетный счет</Label>
                                    <Input className={inputClass} required value={formData.jur_face?.checking_account || ''} onChange={e => handleFaceChange('jur_face', 'checking_account', e.target.value)} />
                                </div>
                                <div className="space-y-2 md:col-span-2">
                                    <Label className={labelClass}>Корреспондентский счет</Label>
                                    <Input className={inputClass} required value={formData.jur_face?.correspondent_account || ''} onChange={e => handleFaceChange('jur_face', 'correspondent_account', e.target.value)} />
                                </div>
                                <div className="space-y-2 md:col-span-2">
                                    <Label className={labelClass}>Юридический адрес</Label>
                                    <Input className={inputClass} required value={formData.jur_face?.jur_address || ''} onChange={e => handleFaceChange('jur_face', 'jur_address', e.target.value)} />
                                </div>
                                <div className="space-y-2 md:col-span-2">
                                    <Label className={labelClass}>Фактический адрес</Label>
                                    <Input className={inputClass} required value={formData.jur_face?.fact_address || ''} onChange={e => handleFaceChange('jur_face', 'fact_address', e.target.value)} />
                                </div>
                            </>
                        )}

                        {formData.org_type === 'ip' && (
                            <>
                                <div className="space-y-2 md:col-span-2">
                                    <Label className={labelClass}>ФИО ИП</Label>
                                    <Input className={inputClass} required value={formData.ip_face?.fio || ''} onChange={e => handleFaceChange('ip_face', 'fio', e.target.value)} />
                                </div>
                                {/* ... Поля для ИП ... */}
                                <div className="space-y-2">
                                    <Label className={labelClass}>ИНН</Label>
                                    <Input className={inputClass} required value={formData.ip_face?.inn || ''} onChange={e => handleFaceChange('ip_face', 'inn', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>ОГРНИП</Label>
                                    <Input className={inputClass} required value={formData.ip_face?.ogrn || ''} onChange={e => handleFaceChange('ip_face', 'ogrn', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>Серия свид-ва</Label>
                                    <Input className={inputClass} type="number" required value={formData.ip_face?.ip_svid_serial || ''} onChange={e => handleFaceChange('ip_face', 'ip_svid_serial', +e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>Номер свид-ва</Label>
                                    <Input className={inputClass} type="number" required value={formData.ip_face?.ip_svid_number || ''} onChange={e => handleFaceChange('ip_face', 'ip_svid_number', +e.target.value)} />
                                </div>
                                <div className="space-y-2 md:col-span-2">
                                    <Label className={labelClass}>Кем выдано</Label>
                                    <Input className={inputClass} required value={formData.ip_face?.ip_svid_givenby || ''} onChange={e => handleFaceChange('ip_face', 'ip_svid_givenby', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>БИК</Label>
                                    <Input className={inputClass} required value={formData.ip_face?.bic || ''} onChange={e => handleFaceChange('ip_face', 'bic', e.target.value)} />
                                </div>
                                <div className="space-y-2">
                                    <Label className={labelClass}>Расчетный счет</Label>
                                    <Input className={inputClass} required value={formData.ip_face?.ras_schot || ''} onChange={e => handleFaceChange('ip_face', 'ras_schot', e.target.value)} />
                                </div>
                                <div className="space-y-2 md:col-span-2">
                                    <Label className={labelClass}>Корр. счет</Label>
                                    <Input className={inputClass} required value={formData.ip_face?.kor_schot || ''} onChange={e => handleFaceChange('ip_face', 'kor_schot', e.target.value)} />
                                </div>
                                <div className="space-y-2 md:col-span-2">
                                    <Label className={labelClass}>Юридический адрес</Label>
                                    <Input className={inputClass} required value={formData.ip_face?.jur_address || ''} onChange={e => handleFaceChange('ip_face', 'jur_address', e.target.value)} />
                                </div>
                            </>
                        )}
                    </div>
                </CardContent>
            </Card>

            <Button type="submit" size="lg" className="w-full bg-[#825e9c] text-white hover:bg-[#6a4c80] py-6 text-lg font-bold uppercase shadow-lg transition-transform hover:-translate-y-1">
                {isEditing ? "Сохранить изменения" : "Создать организацию"}
            </Button>
        </form>
    );
}