import {api, DefaultErrorHandler, Message} from "@/api/api";
import {useUserStore} from "@/context/user-store";

export async function MakeTransfer(
    toType: 'user' | 'org' | 'project',
    toId: number,
    amount: number,
    setMessage: (msg: Message) => void
): Promise<boolean> {
    try {
        const user = useUserStore.getState().user;
        if (!user) {
            setMessage({isError: true, message: "Вы не авторизованы"});
            return false;
        }

        await api.post('/tx/transfer', {
            from_type: 'user',
            from_id: user.id,
            to_type: toType,
            to_id: toId,
            amount: amount
        });

        // После успешного перевода нужно обновить баланс юзера в интерфейсе
        // Для этого дергаем обновление профиля или обновляем стор вручную
        // Лучше перезапросить профиль, но пока просто уменьшим визуально:
        useUserStore.setState({
            user: {...user, balance: user.balance - amount}
        });

        setMessage({isError: false, message: "Перевод выполнен успешно!"});
        return true;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return false;
    }
}

export interface Transaction {
    id: number;
    from_id: number;
    to_id: number;
    amount: number;
    type: string; // 'user_to_project', 'project_to_user', etc.
    created_at: string;
    // Опционально, если бэк обогащает данными, но пока используем ID
}

export async function GetTransactions(filters: {
    user_id?: number,
    project_id?: number,
    limit?: number
}): Promise<Transaction[]> {
    try {
        // Формируем query params для Transaction Service
        // Предполагаем, что сервис поддерживает фильтрацию.
        // Если нет - придется фильтровать на клиенте, но сделаем вид, что API умный.
        // В твоем описании эндпоинтов был /get с параметрами.

        /*
           ВНИМАНИЕ: Так как в Transaction Service мы реализовали пока только /create (Invest/Transfer),
           для истории нам нужно либо добавить GET эндпоинт на бэке, либо (если он есть в задании, но мы пропустили) использовать его.

           Если на бэке нет GET /transactions, давай предположим, что мы его добавили в Repo/Handler транзакций.
           (Я добавлю мок-запрос, который сработает, если ты реализуешь GET /api/tx/history на бэке).
        */

        // Пока вернем пустой массив, чтобы код компилировался,
        // но ниже я напишу код для UserView, который будет готов к данным.

        // const params = new URLSearchParams();
        // if (filters.user_id) params.append("from_id", filters.user_id.toString());
        // ...
        // const res = await api.get(`/tx/history?${params.toString()}`);
        // return res.data;

        return [];
    } catch (e) {
        console.warn(e);
        return [];
    }
}