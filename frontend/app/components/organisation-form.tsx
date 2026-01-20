'use client'

import React, {useEffect, useState} from 'react';
import {Organisation, OrgType} from "@/api/organisation";
import {Input} from "@/app/components/ui/input";
import {Button} from "@/app/components/ui/button";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/app/components/ui/select";
import {Card, CardContent, CardHeader, CardTitle} from "@/app/components/ui/card";
import {Label} from "@/app/components/ui/label";

interface OrganisationFormProps {
    initialData?: Partial<Organisation>;
    onSubmit: (data: Partial<Organisation>) => void;
    isEditing?: boolean;
}

export default function OrganisationForm({initialData, onSubmit, isEditing = false}: OrganisationFormProps) {
    const [formData, setFormData] = useState<Partial<Organisation>>({
        org_type: 'phys',
        phys_face: {} as any,
        jur_face: {} as any,
        ip_face: {} as any,
        ...initialData
    });

    const handleChange = (field: string, value: any) => {
        setFormData(prev => ({...prev, [field]: value}));
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
        // Очищаем лишние поля перед отправкой
        const cleanData = {...formData};
        if (cleanData.org_type === 'phys') { delete cleanData.jur_face; delete cleanData.ip_face; }
        if (cleanData.org_type === 'jur') { delete cleanData.phys_face; delete cleanData.ip_face; }
        if (cleanData.org_type === 'ip') { delete cleanData.phys_face; delete cleanData.jur_face; }

        onSubmit(cleanData);
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-6">
            <Card>
                <CardHeader>
                    <CardTitle>Основные данные</CardTitle>
                </CardHeader>
                <CardContent className="grid gap-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <Label>Название (для отображения)</Label>
                            <Input
                                required
                                value={formData.name || ''}
                                onChange={e => handleChange('name', e.target.value)}
                                placeholder='Например: "Мой Стартап"'
                            />
                        </div>
                        <div className="space-y-2">
                            <Label>Email организации</Label>
                            <Input
                                required type="email"
                                value={formData.email || ''}
                                onChange={e => handleChange('email', e.target.value)}
                            />
                        </div>
                    </div>

                    {!isEditing && (
                        <div className="space-y-2">
                            <Label>Тип юридического лица</Label>
                            <Select
                                value={formData.org_type}
                                onValueChange={(val: OrgType) => handleChange('org_type', val)}
                            >
                                <SelectTrigger>
                                    <SelectValue placeholder="Выберите тип" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="phys">Физическое лицо</SelectItem>
                                    <SelectItem value="jur">Юридическое лицо (ООО, АО)</SelectItem>
                                    <SelectItem value="ip">Индивидуальный предприниматель</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>
                    )}
                </CardContent>
            </Card>

            {/* Рендеринг полей в зависимости от типа */}
            <Card>
                <CardHeader>
                    <CardTitle>Реквизиты</CardTitle>
                </CardHeader>
                <CardContent className="grid gap-4 grid-cols-1 md:grid-cols-2">

                    {formData.org_type === 'phys' && (
                        <>
                            <Input placeholder="ФИО" required
                                   value={formData.phys_face?.fio || ''}
                                   onChange={e => handleFaceChange('phys_face', 'fio', e.target.value)} />
                            <Input placeholder="ИНН" required
                                   value={formData.phys_face?.inn || ''}
                                   onChange={e => handleFaceChange('phys_face', 'inn', e.target.value)} />
                            <Input placeholder="Серия паспорта" type="number" required
                                   value={formData.phys_face?.passport_series || ''}
                                   onChange={e => handleFaceChange('phys_face', 'passport_series', +e.target.value)} />
                            <Input placeholder="Номер паспорта" type="number" required
                                   value={formData.phys_face?.passport_number || ''}
                                   onChange={e => handleFaceChange('phys_face', 'passport_number', +e.target.value)} />
                            <Input placeholder="Кем выдан" required className="md:col-span-2"
                                   value={formData.phys_face?.passport_givenby || ''}
                                   onChange={e => handleFaceChange('phys_face', 'passport_givenby', e.target.value)} />
                            <Input placeholder="БИК Банка" required
                                   value={formData.phys_face?.bic || ''}
                                   onChange={e => handleFaceChange('phys_face', 'bic', e.target.value)} />
                            <Input placeholder="Расчетный счет" required
                                   value={formData.phys_face?.checking_account || ''}
                                   onChange={e => handleFaceChange('phys_face', 'checking_account', e.target.value)} />
                            <Input placeholder="Корр. счет" required
                                   value={formData.phys_face?.correspondent_account || ''}
                                   onChange={e => handleFaceChange('phys_face', 'correspondent_account', e.target.value)} />
                            <Input placeholder="Адрес регистрации" required className="md:col-span-2"
                                   value={formData.phys_face?.registration_address || ''}
                                   onChange={e => handleFaceChange('phys_face', 'registration_address', e.target.value)} />
                            <Input placeholder="Почтовый адрес" required className="md:col-span-2"
                                   value={formData.phys_face?.post_address || ''}
                                   onChange={e => handleFaceChange('phys_face', 'post_address', e.target.value)} />
                        </>
                    )}

                    {formData.org_type === 'jur' && (
                        <>
                            <Input placeholder="Полное наименование" required className="md:col-span-2"
                                   value={formData.jur_face?.full_organisation_name || ''}
                                   onChange={e => handleFaceChange('jur_face', 'full_organisation_name', e.target.value)} />
                            <Input placeholder="Краткое наименование" required
                                   value={formData.jur_face?.short_organisation_name || ''}
                                   onChange={e => handleFaceChange('jur_face', 'short_organisation_name', e.target.value)} />
                            <Input placeholder="ИНН" required
                                   value={formData.jur_face?.inn || ''}
                                   onChange={e => handleFaceChange('jur_face', 'inn', e.target.value)} />
                            <Input placeholder="ОГРН" required
                                   value={formData.jur_face?.ogrn || ''}
                                   onChange={e => handleFaceChange('jur_face', 'ogrn', e.target.value)} />
                            <Input placeholder="КПП" required
                                   value={formData.jur_face?.kpp || ''}
                                   onChange={e => handleFaceChange('jur_face', 'kpp', e.target.value)} />
                            <Input placeholder="Должность руководителя" required
                                   value={formData.jur_face?.position || ''}
                                   onChange={e => handleFaceChange('jur_face', 'position', e.target.value)} />
                            <Input placeholder="Действует на основании" required
                                   value={formData.jur_face?.acts_on_base || ''}
                                   onChange={e => handleFaceChange('jur_face', 'acts_on_base', e.target.value)} />

                            <div className="md:col-span-2 border-t my-2"></div>

                            <Input placeholder="БИК Банка" required
                                   value={formData.jur_face?.bic || ''}
                                   onChange={e => handleFaceChange('jur_face', 'bic', e.target.value)} />
                            <Input placeholder="Расчетный счет" required
                                   value={formData.jur_face?.checking_account || ''}
                                   onChange={e => handleFaceChange('jur_face', 'checking_account', e.target.value)} />
                            <Input placeholder="Корр. счет" required
                                   value={formData.jur_face?.correspondent_account || ''}
                                   onChange={e => handleFaceChange('jur_face', 'correspondent_account', e.target.value)} />
                            <Input placeholder="Юр. Адрес" required className="md:col-span-2"
                                   value={formData.jur_face?.jur_address || ''}
                                   onChange={e => handleFaceChange('jur_face', 'jur_address', e.target.value)} />
                            <Input placeholder="Факт. Адрес" required className="md:col-span-2"
                                   value={formData.jur_face?.fact_address || ''}
                                   onChange={e => handleFaceChange('jur_face', 'fact_address', e.target.value)} />
                        </>
                    )}

                    {formData.org_type === 'ip' && (
                        <>
                            <Input placeholder="ФИО ИП" required
                                   value={formData.ip_face?.fio || ''}
                                   onChange={e => handleFaceChange('ip_face', 'fio', e.target.value)} />
                            <Input placeholder="ИНН" required
                                   value={formData.ip_face?.inn || ''}
                                   onChange={e => handleFaceChange('ip_face', 'inn', e.target.value)} />
                            <Input placeholder="ОГРНИП" required
                                   value={formData.ip_face?.ogrn || ''}
                                   onChange={e => handleFaceChange('ip_face', 'ogrn', e.target.value)} />
                            <Input placeholder="Серия свид-ва" type="number" required
                                   value={formData.ip_face?.ip_svid_serial || ''}
                                   onChange={e => handleFaceChange('ip_face', 'ip_svid_serial', +e.target.value)} />
                            <Input placeholder="Номер свид-ва" type="number" required
                                   value={formData.ip_face?.ip_svid_number || ''}
                                   onChange={e => handleFaceChange('ip_face', 'ip_svid_number', +e.target.value)} />
                            <Input placeholder="Кем выдано" required className="md:col-span-2"
                                   value={formData.ip_face?.ip_svid_givenby || ''}
                                   onChange={e => handleFaceChange('ip_face', 'ip_svid_givenby', e.target.value)} />

                            <div className="md:col-span-2 border-t my-2"></div>

                            <Input placeholder="БИК Банка" required
                                   value={formData.ip_face?.bic || ''}
                                   onChange={e => handleFaceChange('ip_face', 'bic', e.target.value)} />
                            <Input placeholder="Расчетный счет" required
                                   value={formData.ip_face?.ras_schot || ''}
                                   onChange={e => handleFaceChange('ip_face', 'ras_schot', e.target.value)} />
                            <Input placeholder="Корр. счет" required
                                   value={formData.ip_face?.kor_schot || ''}
                                   onChange={e => handleFaceChange('ip_face', 'kor_schot', e.target.value)} />
                            <Input placeholder="Юр. Адрес" required className="md:col-span-2"
                                   value={formData.ip_face?.jur_address || ''}
                                   onChange={e => handleFaceChange('ip_face', 'jur_address', e.target.value)} />
                        </>
                    )}

                </CardContent>
            </Card>

            <Button type="submit" size="lg" className="w-full bg-light-green text-black hover:bg-green-400">
                {isEditing ? "Сохранить изменения" : "Создать организацию"}
            </Button>
        </form>
    );
}