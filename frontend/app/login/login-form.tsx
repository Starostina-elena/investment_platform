import styles from './registration-form.module.css';
import {useEffect, useState} from 'react';
import Spinner from "@/app/components/spinner";
import MessageComponent from "@/app/components/message";
import {useRouter} from "next/navigation";
import {useUserStore} from "@/context/user-store";
import {Message} from "@/api/api";
import {Login} from "@/api/auth";

interface FormData {
    email: string;
    password: string;
}

export default function LoginForm({onRegisterClick}: { onRegisterClick?: () => void }) {
    const [formData, setFormData] = useState<FormData>({
        email: '',
        password: ''
    });

    const [response, setResponse] = useState<Message | null>(null)
    const [requestSent, setRequestSent] = useState<boolean>(false);
    const router = useRouter();
    const user = useUserStore((state) => state);

    useEffect(() => {
        response && setRequestSent(false);
        if (user.user)
            router.push('/');
    }, [response, router, user]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        Login(formData.email, formData.password, setResponse);
    };

    return (
        <form className={styles.login_form} onSubmit={handleSubmit}>
            <div className={styles.auth_options}>
                <h2 onClick={onRegisterClick}>Регистрация</h2>
                <h2>Войти</h2>
            </div>
            <hr style={{width: 'calc(100% + 40px)', margin: '10px -20px'}}/>
            <input
                className={styles.input_field}
                placeholder="Email"
                value={formData.email}
                onChange={(e) => setFormData({...formData, email: e.target.value})}
            />
            <div className={styles.password_field}>
                <input
                    className={styles.input_field}
                    placeholder="Пароль"
                    type="password"
                    value={formData.password}
                    onChange={(e) => setFormData({...formData, password: e.target.value})}
                />
                <a href="#" className={styles.forgot_link}>Забыли?</a>
            </div>
            <MessageComponent message={response}/>
            <button disabled={requestSent} className={styles.submit_button}>{requestSent &&
                <Spinner size={30} style={{margin: "-11px 0 -11px -32px", paddingRight: "32px"}}/>}Войти
            </button>
            <hr style={{width: 'calc(100% + 40px)', margin: '10px -20px'}}/>
            <p className={styles.agreement_text}>
                Нажимая кнопку «Войти» вы принимаете условия пользовательского соглашения и
                даете согласие на обработку персональных данных
            </p>
        </form>
    );
}