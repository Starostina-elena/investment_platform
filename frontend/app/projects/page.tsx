'use client'

import React, {useEffect, useState, Suspense} from 'react';
import {useSearchParams, useRouter} from 'next/navigation';
import {GetProjects, Project} from "@/api/project";
import ProjectPreviewNew from "@/app/components/project-preview-new";
import Spinner from "@/app/components/spinner";
import {CATEGORIES} from "@/app/globals";
import styles from "./projects.module.css"; // Свои стили для каталога

export default function ProjectsPage() {
    return (
        <Suspense fallback={<div style={{display: 'flex', justifyContent: 'center', padding: '10rem'}}><Spinner/></div>}>
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
                // Временная фильтрация на фронте, если бэкенд отдает все подряд
                let filtered = data;
                if (searchQuery) {
                    const lowerQ = searchQuery.toLowerCase();
                    filtered = filtered.filter(p => p.name.toLowerCase().includes(lowerQ) || p.quick_peek.toLowerCase().includes(lowerQ));
                }
                if (selectedCategory) {
                    filtered = filtered.filter(p => p.category === selectedCategory);
                }
                setProjects(filtered);
            })
            .finally(() => setLoading(false));
    }, [searchParams]);

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
        <div className={styles.container}>
            <h1 className={styles.title}>Каталог проектов</h1>

            {/* Фильтры и поиск */}
            <div className={styles.filters_container}>

                {/* Строка поиска */}
                <form onSubmit={handleSearch} className={styles.search_form}>
                    <input
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        placeholder="Найти проект..."
                        className={styles.search_input}
                    />
                    <button type="submit" className={styles.search_button}>
                        Найти
                    </button>
                </form>

                {/* Категории */}
                <div className={styles.categories_list}>
                    <button
                        className={`${styles.category_btn} ${selectedCategory === '' ? styles.category_btn_active : ''}`}
                        onClick={() => handleCategoryClick('')}
                    >
                        Все
                    </button>
                    {Object.keys(CATEGORIES).map(cat => (
                        <button
                            key={cat}
                            className={`${styles.category_btn} ${selectedCategory === cat ? styles.category_btn_active : ''}`}
                            onClick={() => handleCategoryClick(cat)}
                        >
                            {cat}
                        </button>
                    ))}
                </div>
            </div>

            {/* Сетка проектов */}
            {loading ? (
                <div style={{display: 'flex', justifyContent: 'center', padding: '5rem'}}>
                    <Spinner />
                </div>
            ) : projects.length > 0 ? (
                <div className={styles.grid_container}>
                    {projects.map(project => (
                        <ProjectPreviewNew key={project.id} project={project} />
                    ))}
                </div>
            ) : (
                <div className={styles.empty_state}>
                    Ничего не найдено по вашему запросу 😔
                </div>
            )}
        </div>
    );
}