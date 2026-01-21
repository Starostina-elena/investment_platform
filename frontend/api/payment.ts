import {api, DefaultErrorHandler, Message} from "@/api/api";
import {useUserStore} from "@/context/user-store";

export async function InitPayment(
    amount: number,
    returnUrl: string,
    setMessage: (msg: Message) => void
): Promise<string | null> {
    try {
        // 1. Берем юзера из стейта, а не из localStorage
        const user = useUserStore.getState().user;

        if (!user) {
            setMessage({isError: true, message: "Ошибка: пользователь не авторизован"});
            return null;
        }

        console.log("Sending payment init:", { user_id: user.id, amount, return_url: returnUrl });

        // 2. Отправляем запрос
        const res = await api.post("/payment/pay/init", {
            entity_type: "user",
            entity_id: user.id,
            amount: amount,
            return_url: returnUrl
        });

        // 3. Возвращаем ссылку
        return res.data.confirmation_url;
    } catch (e: any) {
        console.error("Payment init error:", e);
        DefaultErrorHandler(setMessage)(e);
        return null;
    }
}

export async function InitWithdraw(
    amount: number,
    payoutToken: string,
    setMessage: (msg: Message) => void
): Promise<string | null> {
    try {
        const user = useUserStore.getState().user;

        if (!user) {
            setMessage({isError: true, message: "Ошибка: пользователь не авторизован"});
            return null;
        }

        console.log("Sending withdraw init:", { user_id: user.id, amount, payout_destination: payoutToken });

        const res = await api.post("/payment/withdraw/init", {
            entity_type: "user",
            entity_id: user.id,
            amount: amount,
            payout_destination: payoutToken
        });

        return res.data.id;
    } catch (e: any) {
        console.error("Withdraw init error:", e);
        DefaultErrorHandler(setMessage)(e);
        return null;
    }
}