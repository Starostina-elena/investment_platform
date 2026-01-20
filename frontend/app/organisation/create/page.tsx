'use client'

import React, {useState} from 'react';
import {useRouter} from 'next/navigation';
import OrganisationForm from "@/app/components/organisation-form";
import {CreateOrganisation} from "@/api/organisation";
import MessageComponent from "@/app/components/message";
import {Message} from "@/api/api";
import styles from "@/app/user-profile/page.module.css"; // Reuse styling

export default function CreateOrganisationPage() {
    const router = useRouter();
    const [message, setMessage] = useState<Message | null>(null);

    const handleCreate = async (data: any) => {
        const newOrg = await CreateOrganisation(data, setMessage);
        if (newOrg) {
            // Редирект на страницу созданной организации
            router.push(`/organisation/${newOrg.id}`);
        }
    };

    return (
        <div className={styles.container} style={{padding: '3rem'}}>
            <div className="max-w-4xl mx-auto">
                <h1 className="text-3xl font-bold text-white mb-6">Регистрация организации</h1>
                <p className="text-gray-300 mb-8">
                    Создайте юридическое лицо для запуска проектов и получения финансирования.
                </p>

                <OrganisationForm onSubmit={handleCreate} />

                <div className="mt-4">
                    <MessageComponent message={message} />
                </div>
            </div>
        </div>
    );
}