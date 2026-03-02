import { useState, useEffect } from 'react';
import { Download } from 'lucide-react';

interface Subject {
  code: string;
  alias?: string;
}

interface TimetableViewProps {
  batch: string;
  subjects: Subject[];
}

export function TimetableView({ batch, subjects }: TimetableViewProps) {
  const [timetableImage, setTimetableImage] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Clean up the object URL to prevent memory leaks
  useEffect(() => {
    return () => {
      if (timetableImage) {
        URL.revokeObjectURL(timetableImage);
      }
    };
  }, [timetableImage]);

  const generateTimetable = async () => {
    if (!batch || subjects.length === 0) {
      setError('Please select a batch and at least one subject');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch('/api/timetable/generate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          batch,
          subjects,
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to generate timetable');
      }

      const blob = await response.blob();
      const imageUrl = URL.createObjectURL(blob);
      setTimetableImage(imageUrl);
    } catch (err) {
      console.error('Error generating timetable:', err);
      setError('Failed to generate timetable. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const downloadTimetable = () => {
    if (!timetableImage) return;

    const link = document.createElement('a');
    link.href = timetableImage;
    link.download = `timetable-${batch}-${new Date().toISOString().split('T')[0]}.png`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  return (
    <div className="w-full max-w-2xl space-y-6">
      <button
        onClick={generateTimetable}
        disabled={!batch || subjects.length === 0 || loading}
        className="w-full px-8 py-5 text-xl font-semibold bg-[#2B2B2B] dark:bg-[#F5F1E8] text-[#F5F1E8] dark:text-[#2B2B2B] rounded-xl hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
      >
        {loading ? 'Generating...' : 'Generate Timetable'}
      </button>

      {error && (
        <div className="p-4 bg-red-100 dark:bg-red-900/30 border border-red-300 dark:border-red-700 rounded-xl text-red-800 dark:text-red-200">
          {error}
        </div>
      )}

      {timetableImage && (
        <div
          className="border border-[#2B2B2B]/15 dark:border-[#F5F1E8]/15 rounded-xl overflow-hidden"
          style={{ aspectRatio: '5120 / 3328' }}
        >
          <img
            src={timetableImage}
            alt="Generated Timetable"
            className="w-full h-full object-cover"
          />
        </div>
      )}

      {!timetableImage && !error && (
        <div
          className="relative border-2 border-dashed border-[#2B2B2B]/20 dark:border-[#F5F1E8]/20 rounded-xl flex items-center justify-center"
          style={{ aspectRatio: '5120 / 3328' }}
        >
          <p className="text-xl text-muted-foreground">
            {loading ? 'Generating your timetable...' : 'Your timetable will appear here'}
          </p>
        </div>
      )}

      <button
        onClick={downloadTimetable}
        disabled={!timetableImage}
        className="w-full px-8 py-5 text-xl font-semibold bg-[#2B2B2B] dark:bg-[#F5F1E8] text-[#F5F1E8] dark:text-[#2B2B2B] rounded-xl hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center justify-center gap-2"
      >
        <Download className="w-5 h-5" />
        Download Timetable
      </button>
    </div>
  );
}