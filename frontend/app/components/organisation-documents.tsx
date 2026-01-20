'use client'

import React, { useState } from 'react';
import {
    DeleteOrgDocument,
    DocType,
    GetOrgDocument,
    UploadOrgDocument
} from "@/api/organisation";
import { Button } from "@/app/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/app/components/ui/card";
import { Upload, Download, Trash2, FileText } from "lucide-react";
import MessageComponent from "@/app/components/message";
import { Message } from "@/api/api";

interface OrganisationDocumentsProps {
    orgId: number;
    orgType: 'phys' | 'jur' | 'ip';
}

const DOC_TYPES: Record<string, { types: DocType[]; labels: Record<DocType, string> }> = {
    phys: {
        types: ['phys_passport_photo', 'phys_passport_propiska', 'phys_svid_uchet'],
        labels: {
            'phys_passport_photo': 'Фото паспорта',
            'phys_passport_propiska': 'Прописка из паспорта',
            'phys_svid_uchet': 'Свидетельство о постановке на налоговый учёт',
            'jur_reg_svid': '',
            'jur_svid_uchet': '',
            'jur_appointment_protocol': '',
            'jur_usn': '',
            'jur_ustav': '',
            'ip_svid_uchet': '',
            'ip_passport_photo': '',
            'ip_passport_propiska': '',
            'ip_usn': '',
            'ip_ogrnip': ''
        }
    },
    jur: {
        types: ['jur_reg_svid', 'jur_svid_uchet', 'jur_appointment_protocol', 'jur_usn', 'jur_ustav'],
        labels: {
            'jur_reg_svid': 'Свидетельство о регистрации юр. лица',
            'jur_svid_uchet': 'Свидетельство о постановке на налоговый учёт',
            'jur_appointment_protocol': 'Протокол о назначении лица',
            'jur_usn': 'Уведомление об УСН',
            'jur_ustav': 'Устав',
            'phys_passport_photo': '',
            'phys_passport_propiska': '',
            'phys_svid_uchet': '',
            'ip_svid_uchet': '',
            'ip_passport_photo': '',
            'ip_passport_propiska': '',
            'ip_usn': '',
            'ip_ogrnip': ''
        }
    },
    ip: {
        types: ['ip_svid_uchet', 'ip_passport_photo', 'ip_passport_propiska', 'ip_usn', 'ip_ogrnip'],
        labels: {
            'ip_svid_uchet': 'Свидетельство о постановке на налоговый учёт',
            'ip_passport_photo': 'Фото паспорта',
            'ip_passport_propiska': 'Прописка из паспорта',
            'ip_usn': 'Уведомление об УСН',
            'ip_ogrnip': 'ОГРНИП',
            'phys_passport_photo': '',
            'phys_passport_propiska': '',
            'phys_svid_uchet': '',
            'jur_reg_svid': '',
            'jur_svid_uchet': '',
            'jur_appointment_protocol': '',
            'jur_usn': '',
            'jur_ustav': ''
        }
    }
};

export default function OrganisationDocuments({ orgId, orgType }: OrganisationDocumentsProps) {
    const [message, setMessage] = useState<Message | null>(null);
    const [uploadingDoc, setUploadingDoc] = useState<DocType | null>(null);

    const docConfig = DOC_TYPES[orgType];
    if (!docConfig) return null;

    const handleUpload = async (docType: DocType, file: File) => {
        setUploadingDoc(docType);
        const success = await UploadOrgDocument(orgId, docType, file, setMessage);
        setUploadingDoc(null);
        if (success) {
            // Документ загружен, можно обновить UI если нужно
        }
    };

    const handleDownload = async (docType: DocType) => {
        const result = await GetOrgDocument(orgId, docType);
        if (result) {
            const url = URL.createObjectURL(result.blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = result.filename;
            a.click();
            URL.revokeObjectURL(url);
        } else {
            setMessage({ isError: true, message: "Документ не найден" });
        }
    };

    const handleDelete = async (docType: DocType) => {
        if (!confirm("Удалить документ?")) return;
        await DeleteOrgDocument(orgId, docType, setMessage);
    };

    return (
        <div className="space-y-4">
            <h3 className="text-xl font-bold text-white">Документы</h3>

            <MessageComponent message={message} />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {docConfig.types.map((docType) => (
                    <Card key={docType} className="bg-[#555652] border-[#666]">
                        <CardHeader>
                            <CardTitle className="text-white text-base flex items-center gap-2">
                                <FileText size={18} />
                                {docConfig.labels[docType]}
                            </CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-2">
                            <div className="flex gap-2">
                                <Button
                                    size="sm"
                                    variant="outline"
                                    className="flex-1 border-[#825e9c] text-[#825e9c]"
                                    disabled={uploadingDoc === docType}
                                    onClick={() => {
                                        const input = document.createElement('input');
                                        input.type = 'file';
                                        input.accept = '.pdf,.jpg,.jpeg,.png,.tiff,.doc,.docx,.xls,.xlsx';
                                        input.onchange = (e) => {
                                            const file = (e.target as HTMLInputElement).files?.[0];
                                            if (file) handleUpload(docType, file);
                                        };
                                        input.click();
                                    }}
                                >
                                    <Upload size={14} className="mr-2" />
                                    {uploadingDoc === docType ? 'Загрузка...' : 'Загрузить'}
                                </Button>
                                <Button
                                    size="sm"
                                    variant="outline"
                                    className="border-blue-500 text-blue-500"
                                    onClick={() => handleDownload(docType)}
                                >
                                    <Download size={14} />
                                </Button>
                                <Button
                                    size="sm"
                                    variant="destructive"
                                    onClick={() => handleDelete(docType)}
                                >
                                    <Trash2 size={14} />
                                </Button>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
}
