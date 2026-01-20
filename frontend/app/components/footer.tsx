import styles from './footer.module.css'
import Image from "next/image";
import facebook from "@/public/facebook_icon.svg";
import instagram from "@/public/instagram_icon.svg";
import vk from "@/public/vk_icon.svg";
import youtube from "@/public/youtube_icon.svg";
import logo from "@/public/logo.png";

const SOCIALS = [
    {
        name: 'Фейсбук',
        link: 'https://www.facebook.com/fundindex',
        icon: facebook
    },
    {
        name: 'Инстаграм',
        link: 'https://www.instagram.com/fundindex/',
        icon: instagram
    },
    {
        name: 'VK',
        link: 'https://vk.com/fundindex',
        icon: vk
    },
    {
        name: 'Ютуб',
        link: 'https://www.youtube.com/@fundindex',
        icon: youtube
    }
];
const LINKS = [
    {
        name: 'О нас',
        link: '/#about'
    },
    {
        name: 'Как это работает',
        link: '/#features'
    },
    {
        name: 'Каталог проектов',
        link: '/#projects'
    },
    {
        name: 'Рекордсмены',
        link: '/#projects'
    },
    {
        name: 'Отзывы',
        link: '/#reviews'
    }
]
export default function Footer(){
    return (
        <footer className={styles.footer}>
            <div>
                <div className={styles.logo}>
                    <img src={logo.src} alt='Логотип'/>
                </div>
                <p className={styles.subscribe_text}>Подпишитесь на нашу рассылку, чтобы быть в курсе новых проектов</p>
                <div className={styles.subscribe}>
                    <label><input type="email" placeholder=""/></label>
                    <button>Подписаться</button>
                </div>
                <p className={styles.subscribe_policy}>Подписываясь, вы соглашаетесь с нашей Политикой конфиденциальности и даете согласие на получение обновлений.</p>
            </div>
            <div className={styles.footer_column}>
                <h4 className={styles.column_title}>Главная</h4>
                {LINKS.map(link => (
                    <a key={link.name} href={link.link} className={styles.social_link}>{link.name}</a>
                ))}
            </div>
            <div className={styles.footer_column} id='contacts'>
                <h4 className={styles.column_title}>Следите за нами</h4>
                {SOCIALS.map(social => (
                    <a key={social.name} href={social.link} className={styles.social_link}>
                        <Image src={social.icon} alt={social.name} width={30} height={30}/>
                        {social.name}
                    </a>
                ))}
            </div>
            <p className={styles.copyright}>© 2025 FundIndex. Все права защищены.</p>
        </footer>
    )
}