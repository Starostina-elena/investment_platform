'use client'
import styles from './page.module.css'
import RegistrationForm from "@/app/login/registration-form";
import LoginForm from "@/app/login/login-form";
import {useState} from "react";


export default function Page() {
    const [registration, setRegistration] = useState(false);

    return (
        <div className={styles.main}>
            {registration ?
                <RegistrationForm onLoginClick={() => setRegistration(false)}/> :
                <LoginForm onRegisterClick={() => setRegistration(true)}/>}
        </div>
    );
}