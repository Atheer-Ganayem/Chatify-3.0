import React, {
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
  useCallback,
} from "react";
import { toast } from "sonner";

type NotificationContextType = {
  notify: (message: string) => void;
  setValues: (values: { active?: boolean; sound?: boolean }) => void;
  active: boolean;
  soundActive: boolean;
};

const NotificationContext = createContext<NotificationContextType | undefined>(
  undefined
);

const LOCALSTORAGE_ACTIVE_KEY = "notifications-active";
const LOCALSTORAGE_SOUND_KEY = "notifications-sound";

const NotificationProvider = ({ children }: { children: React.ReactNode }) => {
  const [active, setActive] = useState(true);
  const [soundActive, setSoundActive] = useState(true);

  const activeRef = useRef(active);
  const soundActiveRef = useRef(soundActive);

  useEffect(() => {
    const storedActive = localStorage.getItem(LOCALSTORAGE_ACTIVE_KEY);
    const storedSound = localStorage.getItem(LOCALSTORAGE_SOUND_KEY);

    if (storedActive !== null) {
      setActive(storedActive === "true");
    }

    if (storedSound !== null) {
      setSoundActive(storedSound === "true");
    }
  }, []);

  useEffect(() => {
    activeRef.current = active;
  }, [active]);

  useEffect(() => {
    soundActiveRef.current = soundActive;
  }, [soundActive]);

  const notify = useCallback((message: string) => {
    if (activeRef.current) toast(message);
    if (soundActiveRef.current) {
      const ding = new Audio("/ding.mp3");
      ding.play();
    }
  }, []);

  const setValues = ({
    active,
    sound,
  }: {
    active?: boolean;
    sound?: boolean;
  }) => {
    if (active !== undefined) {
      setActive(active);
      localStorage.setItem(LOCALSTORAGE_ACTIVE_KEY, active ? "true" : "false");
    }

    if (sound !== undefined) {
      setSoundActive(sound);
      localStorage.setItem(LOCALSTORAGE_SOUND_KEY, sound ? "true" : "false");
    }
  };

  return (
    <NotificationContext.Provider
      value={{ notify, setValues, active, soundActive }}
    >
      {children}
    </NotificationContext.Provider>
  );
};

export const useNotification = () => {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error(
      "useNotification must be used within a NotificationProvider"
    );
  }
  return context;
};

export default NotificationProvider;
