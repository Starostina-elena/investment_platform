'use client'
import CourseSection from '../components/course-section';
import FAQSection from '../components/faq-section';
import ProjectTypeChooser from "@/app/components/project-type-chooser";
import styles from './page.module.css';

const HomePage = () => {
    return (
        <div className={styles.main}>
            <ProjectTypeChooser/>
            <CourseSection/>
            <FAQSection/>
        </div>

    );
};

export default HomePage;