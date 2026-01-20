'use client'
import styles from './header.module.css';
import {CATEGORIES} from "@/app/globals";
import Link from "next/link";
import {useState} from "react";

const CategoriesDropdown = () => {
    const [category, setCategory] = useState<number>(-1);

    return (
        <div className={styles.dropdown}>
            <div className={styles.dropdown_column}>
                <h4 className={styles.column_title}>Категории проектов</h4>
                {Object.keys(CATEGORIES).map((e, i) => (
                    <p key={e} className={styles.dropdown_link + ' ' + (category === i ? styles.dropdown_link_active : '')}
                       onMouseEnter={() => setCategory(i)}>{e}</p>
                ))}
            </div>
            {category >= 0 && <div className={styles.dropdown_splitter}></div>}
            {category >= 0 && <div className={styles.dropdown_column}>
                <h4 className={styles.column_title}>{Object.keys(CATEGORIES)[category]}</h4>
                {CATEGORIES[Object.keys(CATEGORIES)[category]].map((e, i) => (
                    <Link href={`/projects?category=${e}`} key={e} className={styles.dropdown_link}>{e}</Link>
                ))}
            </div>}
        </div>
    );
};

export default CategoriesDropdown;