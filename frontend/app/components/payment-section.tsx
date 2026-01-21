'use client'
import { useState, useEffect, useRef } from "react";
import { Button } from "@/app/components/ui/button";
import { Input } from "@/app/components/ui/input";
import { Label } from "@/app/components/ui/label";
import MessageComponent from "@/app/components/message";
import { Message } from "@/api/api";
import Spinner from "@/app/components/spinner";
import { InitPayment, InitWithdraw } from "@/api/payment";
import { ChevronDown, ChevronUp, CreditCard } from "lucide-react";

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

export default function PaymentSection() {
    // Top-up state
    const [topupAmount, setTopupAmount] = useState<string>("1000");
    const [topupLoading, setTopupLoading] = useState(false);
    const [topupMessage, setTopupMessage] = useState<Message | null>(null);
    const [topupOpen, setTopupOpen] = useState(false);

    // Withdraw state
    const [withdrawAmount, setWithdrawAmount] = useState<string>("1000");
    const [withdrawLoading, setWithdrawLoading] = useState(false);
    const [withdrawMessage, setWithdrawMessage] = useState<Message | null>(null);
    const [withdrawOpen, setWithdrawOpen] = useState(false);
    const [showWithdrawForm, setShowWithdrawForm] = useState(true);
    const [selectedCard, setSelectedCard] = useState<PayoutData | null>(null);
    const payoutsDataRef = useRef<any>(null);
    const widgetRenderedRef = useRef(false);

    const agentId = process.env.NEXT_PUBLIC_YOOKASSA_AGENT_ID || "default_agent_id";

    // Initialize withdraw widget when section opens
    useEffect(() => {
        if (withdrawOpen && showWithdrawForm && !widgetRenderedRef.current) {
            const script = document.createElement("script");
            script.src = "https://yookassa.ru/payouts-data/3.1.0/widget.js";
            script.async = true;
            script.onload = () => {
                initializeWidget();
            };
            document.body.appendChild(script);

            return () => {
                if (payoutsDataRef.current && typeof payoutsDataRef.current.clearListeners === 'function') {
                    payoutsDataRef.current.clearListeners();
                }
                widgetRenderedRef.current = false;
            };
        }
    }, [withdrawOpen, showWithdrawForm]);

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
                            control_primary: '#D1D5DB',
                            control_primary_content: '#000000',
                            background: '#FFFFFF',
                            text: '#091655',
                            border: '#B7FF00',
                            control_secondary: '#E5E7EB'
                        }
                    }
                });

                payoutsDataRef.current.render('payout-form')
                    .then(() => {
                        widgetRenderedRef.current = true;
                    })
                    .catch((error: any) => {
                        setWithdrawMessage({
                            isError: true,
                            message: "Ошибка загрузки формы. Пожалуйста, перезагрузите страницу."
                        });
                    });
            } catch (error: any) {
                setWithdrawMessage({
                    isError: true,
                    message: "Ошибка инициализации виджета. Пожалуйста, проверьте конфигурацию."
                });
            }
        }
    };

    const handlePayoutSuccess = (data: PayoutData) => {
        setSelectedCard(data);
        setShowWithdrawForm(false);
    };

    const handlePayoutError = (error: string) => {
        const userMessage = ERROR_MESSAGES[error] || `Ошибка: ${error}`;
        setWithdrawMessage({
            isError: true,
            message: userMessage
        });
    };

    const handleTopUp = async () => {
        const value = parseFloat(topupAmount);
        if (isNaN(value) || value <= 0) {
            setTopupMessage({isError: true, message: "Введите корректную сумму"});
            return;
        }

        setTopupLoading(true);
        setTopupMessage(null);

        const returnUrl = window.location.href;
        const url = await InitPayment(value, returnUrl, setTopupMessage);

        if (url) {
            window.location.href = url;
        } else {
            setTopupLoading(false);
        }
    };

    const handleWithdraw = async () => {
        if (!selectedCard) {
            setWithdrawMessage({ isError: true, message: "Ошибка: карта не выбрана" });
            return;
        }

        const value = parseFloat(withdrawAmount);
        if (isNaN(value) || value <= 0) {
            setWithdrawMessage({ isError: true, message: "Введите корректную сумму" });
            return;
        }

        if (value < 100) {
            setWithdrawMessage({ isError: true, message: "Минимальная сумма вывода: 100 ₽" });
            return;
        }

        setWithdrawLoading(true);
        setWithdrawMessage(null);

        const result = await InitWithdraw(value, selectedCard.payout_token, setWithdrawMessage);

        if (result) {
            setWithdrawMessage({
                isError: false,
                message: `Запрос на вывод ${value} ₽ успешно отправлен. ID: ${result}`
            });
            setTimeout(() => {
                setWithdrawOpen(false);
                setShowWithdrawForm(true);
                setSelectedCard(null);
                setWithdrawAmount("1000");
                widgetRenderedRef.current = false;
            }, 2000);
        }

        setWithdrawLoading(false);
    };

    const handleCloseWithdraw = () => {
        if (payoutsDataRef.current && typeof payoutsDataRef.current.clearListeners === 'function') {
            payoutsDataRef.current.clearListeners();
        }
        widgetRenderedRef.current = false;
        setWithdrawOpen(false);
        setShowWithdrawForm(true);
        setSelectedCard(null);
        setWithdrawMessage(null);
    };

    return (
        <div style={{ padding: '20px', margin: '0', width: '100%', display: 'block' }}>
            <div className="border-2 border-gray-200 rounded-lg" style={{ marginBottom: '20px' }}>
                <button
                    onClick={() => setTopupOpen(!topupOpen)}
                    className="w-full px-8 bg-gray-50 hover:bg-gray-100 transition-colors"
                    style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', height: '64px', gap: '16px', minWidth: '180px' }}
                >
                    <span className="font-bold text-lg text-gray-900">Пополнить баланс</span>
                    {topupOpen ? (
                        <ChevronUp className="w-6 h-6 text-gray-600" style={{ flexShrink: 0 }} />
                    ) : (
                        <ChevronDown className="w-6 h-6 text-gray-600" style={{ flexShrink: 0 }} />
                    )}
                </button>
                
                {topupOpen && (
                    <div style={{ padding: '20px', display: 'flex', flexDirection: 'column', gap: '5px' }}>
                        <div className="bg-blue-50 p-10 rounded-lg border border-blue-200">
                            <Label htmlFor="topup-amount" className="text-sm font-bold block mb-4 text-gray-900">
                                Сумма (₽)
                            </Label>
                            <Input
                                id="topup-amount"
                                type="number"
                                value={topupAmount}
                                min={100}
                                onChange={(e) => setTopupAmount(e.target.value)}
                                className="bg-white border-blue-300 text-black font-bold text-lg p-5 w-full h-16"
                            />
                            <p className="text-xs text-blue-700 mt-4">Минимальная сумма: 100 ₽</p>
                        </div>

                        <div className="bg-gradient-to-br from-[#B7FF00] to-[#a6e600] p-6 rounded-lg">
                            <p className="text-sm text-gray-900 mb-2">
                                <strong>Итого к оплате:</strong>
                            </p>
                            <p className="text-4xl font-bold text-gray-900">
                                {parseFloat(topupAmount || "0").toLocaleString()} ₽
                            </p>
                        </div>

                        <div className="bg-gray-50 p-5 rounded-lg border border-gray-200">
                            <p className="text-sm text-gray-700">
                                Вы будете перенаправлены на страницу ЮKassa для оплаты банковской картой.
                            </p>
                        </div>

                        <MessageComponent message={topupMessage} />
                        
                        <Button 
                            onClick={handleTopUp} 
                            disabled={topupLoading}
                            style={{ marginTop: '24px', padding: '16px 24px', minHeight: '56px', maxWidth: '400px', alignSelf: 'flex-start' }}
                            className="bg-[#B7FF00] text-black hover:bg-[#a6e600] font-bold text-base h-auto"
                        >
                            {topupLoading ? <Spinner size={20} /> : "Перейти к оплате"}
                        </Button>
                    </div>
                )}
            </div>

            <div className="border-2 border-gray-200 rounded-lg" id="withdraw-section">
                <button
                    onClick={() => setWithdrawOpen(!withdrawOpen)}
                    className="w-full px-8 bg-gray-50 hover:bg-gray-100 transition-colors"
                    style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', height: '64px', gap: '16px', minWidth: '180px' }}
                >
                    <span className="font-bold text-lg text-gray-900">Вывести деньги</span>
                    {withdrawOpen ? (
                        <ChevronUp className="w-6 h-6 text-gray-600" style={{ flexShrink: 0 }} />
                    ) : (
                        <ChevronDown className="w-6 h-6 text-gray-600" style={{ flexShrink: 0 }} />
                    )}
                </button>
                
                {withdrawOpen && (
                    <div style={{ padding: '20px', display: 'flex', flexDirection: 'column', gap: '5px' }}>
                        {showWithdrawForm && (
                            <>
                                <div className="bg-blue-50 p-6 rounded-lg border border-blue-200">
                                    <p className="text-sm text-gray-700 mb-4">
                                        Ввведите данные вашей карты и выберите "Добавить карту". Мы не сохраняем эти данные.
                                    </p>
                                </div>
                                
                                <div className="w-full">
                                    <div 
                                        id="payout-form" 
                                        style={{ display: 'flex', justifyContent: 'flex-start', width: '100%', textAlign: 'left', margin: '0', marginLeft: '-200px', padding: '12px 12px 12px 20px', flex: '1' }}
                                        className="border-2 border-gray-300 rounded-lg bg-white min-h-[360px] overflow-auto"
                                    ></div>
                                </div>
                                
                                <MessageComponent message={withdrawMessage} />
                            </>
                        )}
                        {!showWithdrawForm && (
                            <>
                                <div className="bg-gradient-to-br from-blue-50 to-blue-100 p-6 rounded-lg border border-blue-300">
                                    <h3 className="font-bold text-base mb-4 text-gray-900">Карта выбрана</h3>
                                    <div className="space-y-2 text-sm">
                                        <div className="text-gray-700">
                                            <strong>Номер карты:</strong> {selectedCard?.first6}***{selectedCard?.last4}
                                        </div>
                                        {selectedCard?.issuer_name && (
                                            <div className="text-gray-700">
                                                <strong>Эмитент:</strong> {selectedCard.issuer_name}
                                            </div>
                                        )}
                                    </div>
                                </div>

                                <div className="bg-blue-50 p-10 rounded-lg border border-blue-200">
                                    <Label htmlFor="withdraw-amount-confirm" className="text-sm font-bold block mb-4 text-gray-900">
                                        Сумма вывода (₽)
                                    </Label>
                                    <Input
                                        id="withdraw-amount-confirm"
                                        type="number"
                                        value={withdrawAmount}
                                        min={100}
                                        onChange={(e) => setWithdrawAmount(e.target.value)}
                                        className="bg-white border-blue-300 text-black font-bold text-lg p-5 w-full h-16"
                                    />
                                    <p className="text-xs text-blue-700 mt-4">Минимальная сумма вывода: 100 ₽</p>
                                </div>

                                <div className="bg-gradient-to-br from-[#B7FF00] to-[#a6e600] p-6 rounded-lg">
                                    <p className="text-sm text-gray-900 mb-2">
                                        <strong>К выплате:</strong>
                                    </p>
                                    <p className="text-4xl font-bold text-gray-900">
                                        {parseFloat(withdrawAmount || "0").toLocaleString()} ₽
                                    </p>
                                </div>

                                <MessageComponent message={withdrawMessage} />

                                <div className="flex" style={{ marginTop: '32px', gap: '20px' }}>
                                    <Button
                                        onClick={() => {
                                            setShowWithdrawForm(true);
                                            setSelectedCard(null);
                                            setWithdrawMessage(null);
                                            if (payoutsDataRef.current && typeof payoutsDataRef.current.clearListeners === 'function') {
                                                payoutsDataRef.current.clearListeners();
                                            }
                                            widgetRenderedRef.current = false;
                                        }}
                                        variant="outline"
                                        disabled={withdrawLoading}
                                        style={{ padding: '16px 24px', minHeight: '56px', flex: '1' }}
                                        className="border-2 border-gray-400 hover:bg-gray-100 text-gray-900 font-bold h-auto"
                                    >
                                        ← Выбрать карту заново
                                    </Button>
                                    <Button
                                        onClick={handleWithdraw}
                                        disabled={withdrawLoading}
                                        style={{ padding: '16px 24px', minHeight: '56px', flex: '1' }}
                                        className="bg-[#B7FF00] text-black hover:bg-[#a6e600] font-bold h-auto"
                                    >
                                        {withdrawLoading ? <Spinner size={20} /> : "Подтвердить вывод"}
                                    </Button>
                                </div>
                            </>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}
