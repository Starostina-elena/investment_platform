'use client'

import React, {useEffect, useState, Suspense} from 'react';
import {useSearchParams, useRouter} from 'next/navigation';
import {GetProjects, Project} from "@/api/project";
import ProjectPreviewNew from "@/app/components/project-preview-new";
import Spinner from "@/app/components/spinner";
import {CATEGORIES} from "@/app/globals";
import styles from "./projects.module.css";
import {Input} from "@/app/components/ui/input";
import {
    Pagination,
    PaginationContent,
    PaginationItem,
    PaginationLink,
    PaginationNext,
    PaginationPrevious,
} from "@/app/components/ui/pagination";

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

    // Переменные состояния из URL
    const searchQuery = searchParams.get('q') || '';
    const selectedCategoryKey = searchParams.get('category') || '';

    // Пагинация
    const [page, setPage] = useState(1);
    const limit = 9;
    const [hasMore, setHasMore] = useState(true);

    // Сброс страницы при изменении фильтров
    useEffect(() => {
        setPage(1);
    }, [searchQuery, selectedCategoryKey]);

    // Загрузка проектов
    useEffect(() => {
        setLoading(true);
        const offset = (page - 1) * limit;

        GetProjects(limit + 1, offset, searchQuery, selectedCategoryKey)
            .then(data => {
                // Проверяем есть ли еще страницы
                if (data.length > limit) {
                    setHasMore(true);
                    setProjects(data.slice(0, limit)); // Берем только limit элементов
                } else {
                    setHasMore(false);
                    setProjects(data);
                }
            })
            .finally(() => setLoading(false));
    }, [searchQuery, selectedCategoryKey, page]);

    const updateUrl = (q: string, cat: string) => {
        const params = new URLSearchParams();
        if (q) params.set('q', q);
        if (cat) params.set('category', cat);
        router.push(`/projects?${params.toString()}`);
    }

    return (
        <div className={styles.container}>
            <h1 className={styles.title}>Каталог проектов</h1>

            <div className={styles.filters_container}>
                {/* Поиск */}
                <div className={styles.search_form}>
                    <Input
                        defaultValue={searchQuery}
                        onChange={(e) => updateUrl(e.target.value, selectedCategoryKey)}
                        placeholder="Найти проект..."
                        className={styles.search_input}
                    />
                </div>

                {/* Фильтры */}
                <div className={styles.categories_list}>
                    <button
                        className={`${styles.category_btn} ${selectedCategoryKey === '' ? styles.category_btn_active : ''}`}
                        onClick={() => updateUrl(searchQuery, '')}
                    >
                        Все
                    </button>
                    {Object.entries(CATEGORIES).map(([key, label]) => (
                        <button
                            key={key}
                            className={`${styles.category_btn} ${selectedCategoryKey === key ? styles.category_btn_active : ''}`}
                            onClick={() => updateUrl(searchQuery, key)}
                        >
                            {label}
                        </button>
                    ))}
                </div>
            </div>

            {/* Сетка */}
            {loading ? (
                <div className="flex justify-center p-20"><Spinner /></div>
            ) : projects.length > 0 ? (
                <>
                    <div className={styles.grid_container}>
                        {projects.map(project => (
                            <ProjectPreviewNew key={project.id} project={project} />
                        ))}
                    </div>

                    {/* Пагинация */}
                    <div className="mt-12">
                        <Pagination>
                            <PaginationContent>
                                <PaginationItem>
                                    <PaginationPrevious
                                        className={page <= 1 ? "pointer-events-none opacity-50" : "cursor-pointer bg-white/10 text-black hover:bg-white/20"}
                                        onClick={() => setPage(p => Math.max(1, p - 1))}
                                    />
                                </PaginationItem>

                                <PaginationItem>
                                    <PaginationLink isActive className="bg-[#DB935B] text-black border-none font-bold">
                                        {page}
                                    </PaginationLink>
                                </PaginationItem>

                                <PaginationItem>
                                    <PaginationNext
                                        className={!hasMore ? "pointer-events-none opacity-50" : "cursor-pointer bg-white/10 text-black hover:bg-white/20"}
                                        onClick={() => setPage(p => p + 1)}
                                    />
                                </PaginationItem>
                            </PaginationContent>
                        </Pagination>
                    </div>
                </>
            ) : (
                <div className={styles.empty_state}>
                    Ничего не найдено 😔
                </div>
            )}
        </div>
    );
}