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
                    <button
                        key={i}
                        type="button"
                        onClick={() => onRouteClick(e)}
                        className={styles.nav_link + ' ' + (currentRoute == e ? styles.nav_link_active : '')}
                    >
                        {i + 1}. {e}
                    </button>
                ))}
            </div>

            <div className={styles.preview_link_container}>
                <button
                    type="button"
                    className={styles.nav_link + ' ' + (currentRoute == 'preview' ? styles.nav_link_active : '')}
                    onClick={onPreviewClick}
                >
                    Предпросмотр
                </button>
            </div>
        </div>
    )
}