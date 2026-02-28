import { useState, useEffect } from 'react';
import { Moon, Sun } from 'lucide-react';
import { TimetableForm } from './components/TimetableForm';
import { TimetableView } from './components/TimetableView';

interface Subject {
  code: string;
  alias?: string;
}

export default function App() {
  const [batch, setBatch] = useState('');
  const [subjects, setSubjects] = useState<Subject[]>([]);
  const [isDark, setIsDark] = useState(false);

  // Initialize theme from system preference
  useEffect(() => {
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    setIsDark(prefersDark);
  }, []);

  // Apply dark mode class
  useEffect(() => {
    if (isDark) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }, [isDark]);

  const toggleTheme = () => {
    setIsDark(!isDark);
  };

  return (
    <div className="min-h-screen bg-background text-foreground transition-colors duration-300">
      {/* Theme Toggle Button */}
      <button
        onClick={toggleTheme}
        className="fixed top-6 right-6 p-3 rounded-full bg-white dark:bg-[#3A3A3A] border border-[#2B2B2B]/15 dark:border-[#F5F1E8]/15 hover:scale-110 transition-all shadow-lg z-50"
        aria-label="Toggle theme"
      >
        {isDark ? <Sun className="w-5 h-5" /> : <Moon className="w-5 h-5" />}
      </button>

      <div className="grid lg:grid-cols-2 gap-8 p-8 lg:p-12 min-h-screen">
        {/* Left Side - Form */}
        <div className="flex flex-col justify-center max-w-xl mx-auto w-full">
          <div className="mb-8">
            <h1 className="mb-2 text-8xl leading-tight font-bold">
              Timetable<br />Builder
            </h1>
          </div>
          
          <TimetableForm
            batch={batch}
            setBatch={setBatch}
            subjects={subjects}
            setSubjects={setSubjects}
          />
        </div>

        {/* Right Side - Timetable */}
        <div className="flex items-center justify-center">
          <TimetableView batch={batch} subjects={subjects} />
        </div>
      </div>
    </div>
  );
}