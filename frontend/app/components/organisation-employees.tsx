'use client'

import React, { useEffect, useState } from 'react';
import {
    AddEmployee,
    DeleteEmployee,
    Employee,
    GetOrganisationEmployees,
    UpdateEmployee
} from "@/api/organisation";
import { Button } from "@/app/components/ui/button";
import { Input } from "@/app/components/ui/input";
import { Label } from "@/app/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/app/components/ui/card";
import { Badge } from "@/app/components/ui/badge";
import { UserPlus, Trash2, Edit2, Save, X } from "lucide-react";
import Spinner from "@/app/components/spinner";
import MessageComponent from "@/app/components/message";
import { Message } from "@/api/api";

interface OrganisationEmployeesProps {
    orgId: number;
}

export default function OrganisationEmployees({ orgId }: OrganisationEmployeesProps) {
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [loading, setLoading] = useState(true);
    const [message, setMessage] = useState<Message | null>(null);
    const [showAddForm, setShowAddForm] = useState(false);
    const [editingId, setEditingId] = useState<number | null>(null);

    const [formData, setFormData] = useState({
        user_id: 0,
        org_account_management: false,
        money_management: false,
        project_management: false
    });

    useEffect(() => {
        loadEmployees();
    }, [orgId]);

    const loadEmployees = async () => {
        setLoading(true);
        const data = await GetOrganisationEmployees(orgId);
        setEmployees(data);
        setLoading(false);
    };

    const handleAdd = async () => {
        if (formData.user_id <= 0) {
            setMessage({ isError: true, message: "Введите корректный ID пользователя" });
            return;
        }

        const success = await AddEmployee(
            orgId,
            formData.user_id,
            {
                org_account_management: formData.org_account_management,
                money_management: formData.money_management,
                project_management: formData.project_management
            },
            setMessage
        );

        if (success) {
            setShowAddForm(false);
            setFormData({
                user_id: 0,
                org_account_management: false,
                money_management: false,
                project_management: false
            });
            loadEmployees();
        }
    };

    const handleUpdate = async (employee: Employee) => {
        const success = await UpdateEmployee(
            orgId,
            employee.user_id,
            {
                org_account_management: employee.org_account_management,
                money_management: employee.money_management,
                project_management: employee.project_management
            },
            setMessage
        );

        if (success) {
            setEditingId(null);
            loadEmployees();
        }
    };

    const handleDelete = async (userId: number) => {
        if (!confirm("Удалить сотрудника?")) return;

        const success = await DeleteEmployee(orgId, userId, setMessage);
        if (success) {
            loadEmployees();
        }
    };

    const togglePermission = (employee: Employee, field: keyof Pick<Employee, 'org_account_management' | 'money_management' | 'project_management'>) => {
        const updated = employees.map(e =>
            e.user_id === employee.user_id ? { ...e, [field]: !e[field] } : e
        );
        setEmployees(updated);
    };

    if (loading) return <div className="flex justify-center p-4"><Spinner /></div>;

    return (
        <div className="space-y-4">
            <div className="flex justify-between items-center">
                <h3 className="text-xl font-bold text-white">Сотрудники</h3>
                <Button onClick={() => setShowAddForm(!showAddForm)} className="bg-[#825e9c]">
                    <UserPlus size={16} className="mr-2" />
                    Добавить сотрудника
                </Button>
            </div>

            <MessageComponent message={message} />

            {showAddForm && (
                <Card className="bg-[#555652] border-[#666]">
                    <CardHeader>
                        <CardTitle className="text-white">Новый сотрудник</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div>
                            <Label className="text-gray-300">ID пользователя</Label>
                            <Input
                                type="number"
                                value={formData.user_id || ''}
                                onChange={(e) => setFormData({ ...formData, user_id: +e.target.value })}
                                className="bg-[#656662] text-white border-[#777]"
                            />
                        </div>

                        <div className="space-y-2">
                            <Label className="text-gray-300">Права доступа:</Label>
                            <div className="flex flex-col gap-2">
                                <label className="flex items-center gap-2 text-white cursor-pointer">
                                    <input
                                        type="checkbox"
                                        checked={formData.org_account_management}
                                        onChange={(e) => setFormData({ ...formData, org_account_management: e.target.checked })}
                                        className="w-4 h-4"
                                    />
                                    Управление аккаунтом организации
                                </label>
                                <label className="flex items-center gap-2 text-white cursor-pointer">
                                    <input
                                        type="checkbox"
                                        checked={formData.money_management}
                                        onChange={(e) => setFormData({ ...formData, money_management: e.target.checked })}
                                        className="w-4 h-4"
                                    />
                                    Управление финансами
                                </label>
                                <label className="flex items-center gap-2 text-white cursor-pointer">
                                    <input
                                        type="checkbox"
                                        checked={formData.project_management}
                                        onChange={(e) => setFormData({ ...formData, project_management: e.target.checked })}
                                        className="w-4 h-4"
                                    />
                                    Управление проектами
                                </label>
                            </div>
                        </div>

                        <div className="flex gap-2">
                            <Button onClick={handleAdd} className="bg-green-600">
                                <Save size={16} className="mr-2" />
                                Добавить
                            </Button>
                            <Button onClick={() => setShowAddForm(false)} variant="outline" className="border-gray-500 text-gray-300">
                                <X size={16} className="mr-2" />
                                Отмена
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            )}

            <div className="space-y-3">
                {employees.length === 0 ? (
                    <Card className="bg-[#555652] border-[#666]">
                        <CardContent className="p-6 text-center text-gray-400">
                            Нет сотрудников
                        </CardContent>
                    </Card>
                ) : (
                    employees.map((employee) => (
                        <Card key={employee.user_id} className="bg-[#555652] border-[#666]">
                            <CardContent className="p-4">
                                <div className="flex justify-between items-start">
                                    <div className="flex-1">
                                        <div className="flex items-center gap-3 mb-2">
                                            <h4 className="text-white font-bold">{employee.nickname}</h4>
                                            <Badge variant="outline" className="text-xs">
                                                ID: {employee.user_id}
                                            </Badge>
                                        </div>
                                        <p className="text-gray-400 text-sm mb-3">{employee.user_email}</p>

                                        <div className="flex flex-wrap gap-2">
                                            {editingId === employee.user_id ? (
                                                <>
                                                    <label className="flex items-center gap-2 text-white text-sm cursor-pointer">
                                                        <input
                                                            type="checkbox"
                                                            checked={employee.org_account_management}
                                                            onChange={() => togglePermission(employee, 'org_account_management')}
                                                            className="w-4 h-4"
                                                        />
                                                        Управление аккаунтом
                                                    </label>
                                                    <label className="flex items-center gap-2 text-white text-sm cursor-pointer">
                                                        <input
                                                            type="checkbox"
                                                            checked={employee.money_management}
                                                            onChange={() => togglePermission(employee, 'money_management')}
                                                            className="w-4 h-4"
                                                        />
                                                        Финансы
                                                    </label>
                                                    <label className="flex items-center gap-2 text-white text-sm cursor-pointer">
                                                        <input
                                                            type="checkbox"
                                                            checked={employee.project_management}
                                                            onChange={() => togglePermission(employee, 'project_management')}
                                                            className="w-4 h-4"
                                                        />
                                                        Проекты
                                                    </label>
                                                </>
                                            ) : (
                                                <>
                                                    {employee.org_account_management && (
                                                        <Badge className="bg-blue-600">Управление аккаунтом</Badge>
                                                    )}
                                                    {employee.money_management && (
                                                        <Badge className="bg-green-600">Финансы</Badge>
                                                    )}
                                                    {employee.project_management && (
                                                        <Badge className="bg-purple-600">Проекты</Badge>
                                                    )}
                                                    {!employee.org_account_management && !employee.money_management && !employee.project_management && (
                                                        <Badge variant="outline" className="text-gray-400">Нет прав</Badge>
                                                    )}
                                                </>
                                            )}
                                        </div>
                                    </div>

                                    <div className="flex gap-2">
                                        {editingId === employee.user_id ? (
                                            <>
                                                <Button
                                                    size="sm"
                                                    onClick={() => handleUpdate(employee)}
                                                    className="bg-green-600"
                                                >
                                                    <Save size={14} />
                                                </Button>
                                                <Button
                                                    size="sm"
                                                    variant="outline"
                                                    onClick={() => {
                                                        setEditingId(null);
                                                        loadEmployees();
                                                    }}
                                                    className="border-gray-500"
                                                >
                                                    <X size={14} />
                                                </Button>
                                            </>
                                        ) : (
                                            <>
                                                <Button
                                                    size="sm"
                                                    variant="outline"
                                                    onClick={() => setEditingId(employee.user_id)}
                                                    className="border-[#825e9c] text-[#825e9c]"
                                                >
                                                    <Edit2 size={14} />
                                                </Button>
                                                <Button
                                                    size="sm"
                                                    variant="destructive"
                                                    onClick={() => handleDelete(employee.user_id)}
                                                >
                                                    <Trash2 size={14} />
                                                </Button>
                                            </>
                                        )}
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    ))
                )}
            </div>
        </div>
    );
}
