'use client'
import { useState, useEffect, useRef } from "react";
import { Button } from "@/app/components/ui/button";
import { Input } from "@/app/components/ui/input";
import { Label } from "@/app/components/ui/label";
import MessageComponent from "@/app/components/message";
import { Message } from "@/api/api";
import Spinner from "@/app/components/spinner";
import { InitWithdraw } from "@/api/payment";
import { CreditCard, ChevronLeft } from "lucide-react";
import SimpleModal from "@/app/components/simple-modal";

declare global {
    interface Window {
        PayoutsData: any;
    }
}

interface PayoutData {
    payout_token: string;
    first6: string;
    last4: string;
    issuer_name: string;
    issuer_country: string;
    card_type: string;
}

const ERROR_MESSAGES: { [key: string]: string } = {
    "card_country_code_error": "Нельзя сделать выплату на банковскую карту, выпущенную в этой стране. Пожалуйста, используйте другую карту.",
    "card_unknown_country_code_error": "Неправильный номер банковской карты: невозможно определить код страны. Пожалуйста, проверьте номер карты.",
    "internal_service_error": "Ошибка обработки карты. Пожалуйста, попробуйте еще раз или используйте другую карту."
};

export default function WithdrawModal() {
    const [amount, setAmount] = useState<string>("1000");
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState<Message | null>(null);
    const [isOpen, setIsOpen] = useState(false);
    const [showForm, setShowForm] = useState(true);
    const [selectedCard, setSelectedCard] = useState<PayoutData | null>(null);
    const payoutsDataRef = useRef<any>(null);
    const widgetRenderedRef = useRef(false);

    // Получить YOOKASSA_AGENT_ID из переменной окружения
    const agentId = process.env.NEXT_PUBLIC_YOOKASSA_AGENT_ID || "default_agent_id";

    useEffect(() => {
        if (isOpen && showForm && !widgetRenderedRef.current) {
            // Загрузить скрипт виджета Юкассы
            const script = document.createElement("script");
            script.src = "https://yookassa.ru/payouts-data/3.1.0/widget.js";
            script.async = true;
            script.onload = () => {
                initializeWidget();
            };
            document.body.appendChild(script);

            return () => {
                // Cleanup при закрытии модального окна
                if (payoutsDataRef.current && typeof payoutsDataRef.current.clearListeners === 'function') {
                    payoutsDataRef.current.clearListeners();
                }
                widgetRenderedRef.current = false;
            };
        }
    }, [isOpen, showForm]);

    const initializeWidget = () => {
        if (typeof window !== 'undefined' && window.PayoutsData) {
            try {
                payoutsDataRef.current = new window.PayoutsData({
                    type: 'payout',
                    account_id: agentId,
                    success_callback: handlePayoutSuccess,
                    error_callback: handlePayoutError,
                    customization: {
                        colors: {
                            control_primary: '#B7FF00',
                            control_primary_content: '#000000',
                            background: '#FFFFFF',
                            text: '#000000',
                            border: '#E0E0E0',
                            control_secondary: '#CCCCCC'
                        }
                    }
                });

                payoutsDataRef.current.render('payout-form')
                    .then(() => {
                        widgetRenderedRef.current = true;
                        console.log("Payout widget rendered successfully");
                    })
                    .catch((error: any) => {
                        console.error("Error rendering payout widget:", error);
                        setMessage({
                            isError: true,
                            message: "Ошибка загрузки формы. Пожалуйста, перезагрузите страницу."
                        });
                    });
            } catch (error: any) {
                console.error("Error initializing payout widget:", error);
                setMessage({
                    isError: true,
                    message: "Ошибка инициализации виджета. Пожалуйста, проверьте конфигурацию."
                });
            }
        }
    };

    const handlePayoutSuccess = (data: PayoutData) => {
        console.log("Payout success:", data);
        setSelectedCard(data);
        setShowForm(false);
    };

    const handlePayoutError = (error: string) => {
        console.error("Payout error:", error);
        const userMessage = ERROR_MESSAGES[error] || `Ошибка: ${error}`;
        setMessage({
            isError: true,
            message: userMessage
        });
    };

    const handleWithdraw = async () => {
        if (!selectedCard) {
            setMessage({ isError: true, message: "Ошибка: карта не выбрана" });
            return;
        }

        const value = parseFloat(amount);
        if (isNaN(value) || value <= 0) {
            setMessage({ isError: true, message: "Введите корректную сумму" });
            return;
        }

        if (value < 100) {
            setMessage({ isError: true, message: "Минимальная сумма вывода: 100 ₽" });
            return;
        }

        setLoading(true);
        setMessage(null);

        const result = await InitWithdraw(value, selectedCard.payout_token, setMessage);

        if (result) {
            setMessage({
                isError: false,
                message: `Запрос на вывод ${value} ₽ успешно отправлен. ID: ${result}`
            });
            setTimeout(() => {
                setIsOpen(false);
                setShowForm(true);
                setSelectedCard(null);
                setAmount("1000");
                widgetRenderedRef.current = false;
            }, 2000);
        }

        setLoading(false);
    };

    const handleClose = () => {
        if (payoutsDataRef.current && typeof payoutsDataRef.current.clearListeners === 'function') {
            payoutsDataRef.current.clearListeners();
        }
        widgetRenderedRef.current = false;
        setIsOpen(false);
        setShowForm(true);
        setSelectedCard(null);
        setMessage(null);
    };

    const handleBackToForm = () => {
        setShowForm(true);
        setSelectedCard(null);
        setMessage(null);
    };

    return (
        <>
            <Button
                onClick={() => setIsOpen(true)}
                variant="outline"
                size="sm"
                className="ml-2 gap-2 border-[#DB935B] text-[#DB935B] hover:bg-[#DB935B] hover:text-black z-10 pointer-events-auto relative"
            >
                <CreditCard className="w-4 h-4" /> Вывести деньги
            </Button>

            <SimpleModal
                isOpen={isOpen}
                onClose={handleClose}
                title={showForm ? "Вывод денежных средств" : "Подтверждение вывода"}
            >
                <div className="space-y-5">
                    {showForm ? (
                        <>
                            <div className="bg-blue-50 p-5 rounded-lg border border-blue-200">
                                <Label htmlFor="withdraw-amount" className="text-sm font-bold block mb-3 text-gray-900">Сумма (₽)</Label>
                                <Input
                                    id="withdraw-amount"
                                    type="number"
                                    value={amount}
                                    min={100}
                                    onChange={(e) => setAmount(e.target.value)}
                                    className="bg-white border-blue-300 text-black font-bold text-lg p-3 w-full"
                                />
                                <p className="text-xs text-blue-700 mt-2">Минимальная сумма вывода: 100 ₽</p>
                            </div>
                            
                            <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                                <p className="text-sm text-gray-700">
                                    Нажмите на форму ниже, введите данные вашей карты и выберите "Далее".
                                </p>
                            </div>
                            
                            <div 
                                id="payout-form" 
                                className="border-2 border-gray-300 rounded-lg p-5 bg-white min-h-[320px] overflow-x-auto"
                            ></div>
                            
                            <MessageComponent message={message} />
                        </>
                    ) : (
                        <>
                            <div className="bg-gradient-to-br from-blue-50 to-blue-100 p-5 rounded-lg border border-blue-300">
                                <h3 className="font-bold text-base mb-4 text-gray-900">Данные банковской карты</h3>
                                <div className="space-y-3 text-sm">
                                    <div className="flex justify-between items-center bg-white p-3 rounded">
                                        <span className="text-gray-700">Номер карты:</span>
                                        <span className="font-mono font-bold text-gray-900">
                                            {selectedCard?.first6}***{selectedCard?.last4}
                                        </span>
                                    </div>
                                    <div className="flex justify-between items-center bg-white p-3 rounded">
                                        <span className="text-gray-700">Система:</span>
                                        <span className="font-bold text-gray-900">{selectedCard?.card_type}</span>
                                    </div>
                                    <div className="flex justify-between items-center bg-white p-3 rounded">
                                        <span className="text-gray-700">Эмитент:</span>
                                        <span className="font-bold text-gray-900">{selectedCard?.issuer_name}</span>
                                    </div>
                                    <div className="flex justify-between items-center bg-white p-3 rounded">
                                        <span className="text-gray-700">Страна:</span>
                                        <span className="font-bold text-gray-900">{selectedCard?.issuer_country}</span>
                                    </div>
                                </div>
                            </div>

                            <div className="bg-gradient-to-br from-[#B7FF00] to-[#a6e600] p-5 rounded-lg">
                                <p className="text-sm text-gray-900 mb-2">
                                    <strong>Сумма вывода:</strong>
                                </p>
                                <p className="text-3xl font-bold text-gray-900">
                                    {parseFloat(amount).toLocaleString()} ₽
                                </p>
                            </div>

                            <MessageComponent message={message} />

                            <div className="flex gap-3">
                                <Button
                                    onClick={handleBackToForm}
                                    variant="outline"
                                    disabled={loading}
                                    className="flex-1 border-2 border-gray-400 hover:bg-gray-100 text-gray-900 font-bold py-6 h-auto"
                                >
                                    <ChevronLeft className="w-4 h-4 mr-2" /> Назад
                                </Button>
                                <Button
                                    onClick={handleWithdraw}
                                    disabled={loading}
                                    className="flex-1 bg-[#B7FF00] text-black hover:bg-[#a6e600] font-bold py-6 h-auto"
                                >
                                    {loading ? <Spinner size={20} /> : "Подтвердить вывод"}
                                </Button>
                            </div>
                        </>
                    )}
                </div>
            </SimpleModal>
        </>
    );
}
