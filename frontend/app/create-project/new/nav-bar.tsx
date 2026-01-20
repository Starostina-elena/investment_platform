import styles from './nav-bar.module.css'

export default function NavBar({routes, onRouteClick, currentRoute, onPreviewClick}: {
    routes: string[],
    onRouteClick: (route: string) => void,
    currentRoute: string,
    onPreviewClick: () => void
}) {


    return (
        <div className={styles.navbar} id='navbar'>
            <div>
                {routes.map((e, i) => (
                    <a href="#navbar" key={i} onClick={onRouteClick.bind(null, e)}
                       className={styles.nav_link + ' ' + (currentRoute == e ? styles.nav_link_active : '')}>
                        {i + 1}. {e}
                    </a>
                ))}
            </div>

            <div className={styles.preview_link_container}>
                <a href="#navbar" className={styles.nav_link + ' ' + (currentRoute == 'preview' ? styles.nav_link_active : '')} onClick={onPreviewClick}>
                    Предпросмотр
                </a>
            </div>
        </div>
    )
}