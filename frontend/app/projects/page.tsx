'use client'

import React, {useEffect, useState, Suspense} from 'react';
import {useSearchParams, useRouter} from 'next/navigation';
import {GetProjects, Project} from "@/api/project";
import ProjectPreviewNew from "@/app/components/project-preview-new";
import Spinner from "@/app/components/spinner";
import {Input} from "@/app/components/ui/input";
import {Button} from "@/app/components/ui/button";
import {CATEGORIES} from "@/app/globals";
import styles from "@/app/page.module.css"; // Переиспользуем стили главной

export default function ProjectsPage() {
    return (
        <Suspense fallback={<div className="flex justify-center p-20"><Spinner/></div>}>
            <Catalog />
        </Suspense>
    );
}

function Catalog() {
    const searchParams = useSearchParams();
    const router = useRouter();

    const [projects, setProjects] = useState<Project[]>([]);
    const [loading, setLoading] = useState(true);

    // Считываем фильтры из URL
    const initialQuery = searchParams.get('q') || '';
    const initialCategory = searchParams.get('category') || '';

    const [searchQuery, setSearchQuery] = useState(initialQuery);
    const [selectedCategory, setSelectedCategory] = useState(initialCategory);

    // Загрузка проектов
    useEffect(() => {
        setLoading(true);
        GetProjects(100, 0, searchQuery, selectedCategory)
            .then(data => {
                // Если бэкенд не умеет фильтровать, делаем это тут (временный фолбек)
                let filtered = data;
                /*
                // Раскомментируй, если бэкенд тупой и возвращает всё подряд
                if (searchQuery) {
                    const lowerQ = searchQuery.toLowerCase();
                    filtered = filtered.filter(p => p.name.toLowerCase().includes(lowerQ) || p.quick_peek.toLowerCase().includes(lowerQ));
                }
                if (selectedCategory) {
                    filtered = filtered.filter(p => p.category === selectedCategory);
                }
                */
                setProjects(filtered);
            })
            .finally(() => setLoading(false));
    }, [searchParams]); // Перезагружаем при изменении URL

    // Обновление URL при поиске
    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        const params = new URLSearchParams();
        if (searchQuery) params.set('q', searchQuery);
        if (selectedCategory) params.set('category', selectedCategory);
        router.push(`/projects?${params.toString()}`);
    };

    // Клики по категориям
    const handleCategoryClick = (cat: string) => {
        const newCat = selectedCategory === cat ? '' : cat;
        setSelectedCategory(newCat);
        const params = new URLSearchParams(searchParams.toString());
        if (newCat) params.set('category', newCat);
        else params.delete('category');
        router.push(`/projects?${params.toString()}`);
    };

    return (
        <div className="min-h-screen bg-[#130622] text-white pt-10 pb-20 px-4 md:px-10">
            <h1 className="text-4xl font-bold mb-8 text-center font-soyuz">Каталог проектов</h1>

            {/* Фильтры и поиск */}
            <div className="max-w-6xl mx-auto mb-12 space-y-6">

                {/* Строка поиска */}
                <form onSubmit={handleSearch} className="flex gap-4">
                    <Input
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        placeholder="Найти проект..."
                        className="bg-[#1e0e31] border-gray-600 text-white text-lg h-12"
                    />
                    <Button type="submit" className="bg-[#825e9c] text-black font-bold h-12 px-8">
                        Найти
                    </Button>
                </form>

                {/* Категории (Основные группы) */}
                <div className="flex flex-wrap gap-2 justify-center">
                    <Button
                        variant={selectedCategory === '' ? "default" : "outline"}
                        className={`rounded-full ${selectedCategory === '' ? 'bg-light-green text-black' : 'text-gray-300 border-gray-600'}`}
                        onClick={() => handleCategoryClick('')}
                    >
                        Все
                    </Button>
                    {Object.keys(CATEGORIES).map(cat => (
                        <Button
                            key={cat}
                            variant={selectedCategory === cat ? "default" : "outline"}
                            className={`rounded-full ${selectedCategory === cat ? 'bg-[#825e9c] text-black border-none' : 'text-gray-300 border-gray-600 hover:bg-white/10'}`}
                            onClick={() => handleCategoryClick(cat)}
                        >
                            {cat}
                        </Button>
                    ))}
                </div>
            </div>

            {/* Сетка проектов */}
            {loading ? (
                <div className="flex justify-center py-20"><Spinner /></div>
            ) : projects.length > 0 ? (
                <div className={styles.projects_container}>
                    {projects.map(project => (
                        <ProjectPreviewNew key={project.id} project={project} />
                    ))}
                </div>
            ) : (
                <div className="text-center py-20 text-gray-500 text-xl">
                    Ничего не найдено по вашему запросу 😔
                </div>
            )}
        </div>
    );
}