import { useState, useRef, useEffect } from 'react';
import { X, ChevronDown, Edit2, Check } from 'lucide-react';

// Maximum character limit for subject aliases
const MAX_ALIAS_LENGTH = 20;

interface Subject {
  code: string;
  alias?: string;
}

interface TimetableFormProps {
  batch: string;
  setBatch: (batch: string) => void;
  subjects: Subject[];
  setSubjects: (subjects: Subject[]) => void;
}

export function TimetableForm({ batch, setBatch, subjects, setSubjects }: TimetableFormProps) {
  const [searchTerm, setSearchTerm] = useState('');
  const [batchSearchTerm, setBatchSearchTerm] = useState('');
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [isBatchDropdownOpen, setIsBatchDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const batchDropdownRef = useRef<HTMLDivElement>(null);

  const [availableSubjects, setAvailableSubjects] = useState<string[]>([]);
  const [batchOptions, setBatchOptions] = useState<string[]>([]);
  const [loadingSubjects, setLoadingSubjects] = useState(true);
  const [loadingBatches, setLoadingBatches] = useState(true);

  const [editingSubjectCode, setEditingSubjectCode] = useState<string | null>(null);
  const [editingAlias, setEditingAlias] = useState('');

  // Fetch batch options from backend (retries every 10s if backend isn't ready)
  useEffect(() => {
    let retryTimeout: ReturnType<typeof setTimeout>;
    let cancelled = false;

    const fetchBatches = async () => {
      try {
        const cached = localStorage.getItem('cached_batches');
        if (cached) {
          setBatchOptions(JSON.parse(cached));
          setLoadingBatches(false);
          return;
        }

        const response = await fetch('/api/timetable/batches');

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        setBatchOptions(data.batches);
        localStorage.setItem('cached_batches', JSON.stringify(data.batches));
        setLoadingBatches(false);
      } catch (error) {
        console.error('Error fetching batches, retrying in 10s...', error);
        if (!cancelled) {
          retryTimeout = setTimeout(fetchBatches, 10000);
        }
      }
    };

    fetchBatches();
    return () => { cancelled = true; clearTimeout(retryTimeout); };
  }, []);

  // Fetch subjects from backend (retries every 10s if backend isn't ready)
  useEffect(() => {
    let retryTimeout: ReturnType<typeof setTimeout>;
    let cancelled = false;

    const fetchSubjects = async () => {
      try {
        const cached = localStorage.getItem('cached_subjects');
        if (cached) {
          setAvailableSubjects(JSON.parse(cached));
          setLoadingSubjects(false);
          return;
        }

        const response = await fetch('/api/timetable/subjects');

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        setAvailableSubjects(data.subjects);
        localStorage.setItem('cached_subjects', JSON.stringify(data.subjects));
        setLoadingSubjects(false);
      } catch (error) {
        console.error('Error fetching subjects, retrying in 10s...', error);
        if (!cancelled) {
          retryTimeout = setTimeout(fetchSubjects, 10000);
        }
      }
    };

    fetchSubjects();
    return () => { cancelled = true; clearTimeout(retryTimeout); };
  }, []);

  const filteredSubjects = availableSubjects.filter(
    (subjectCode) =>
      subjectCode.toLowerCase().includes(searchTerm.toLowerCase()) &&
      !subjects.some(s => s.code === subjectCode)
  );

  const addSubject = (subjectCode: string) => {
    if (!subjects.some(s => s.code === subjectCode)) {
      setSubjects([...subjects, { code: subjectCode }]);
      setSearchTerm('');
      setIsDropdownOpen(false);
    }
  };

  const removeSubject = (subjectCode: string) => {
    setSubjects(subjects.filter((subject) => subject.code !== subjectCode));
  };

  const startEditingAlias = (subjectCode: string, currentAlias?: string) => {
    setEditingSubjectCode(subjectCode);
    setEditingAlias(currentAlias || '');
  };

  const saveAlias = (subjectCode: string) => {
    const trimmedAlias = editingAlias.trim();
    setSubjects(
      subjects.map((subject) =>
        subject.code === subjectCode
          ? { ...subject, alias: trimmedAlias || undefined }
          : subject
      )
    );
    setEditingSubjectCode(null);
    setEditingAlias('');
  };

  const filteredBatches = batchOptions.filter(
    (b) => b.toLowerCase().includes(batchSearchTerm.toLowerCase())
  );

  const selectBatch = (selectedBatch: string) => {
    setBatch(selectedBatch);
    setBatchSearchTerm('');
    setIsBatchDropdownOpen(false);
  };

  // Close dropdowns when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsDropdownOpen(false);
      }
      if (batchDropdownRef.current && !batchDropdownRef.current.contains(event.target as Node)) {
        setIsBatchDropdownOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <div className="space-y-10">
      {/* Batch Dropdown */}
      <div className="space-y-4">
        <label htmlFor="batch-select" className="block text-2xl">
          Batch
        </label>
        <div className="relative" ref={batchDropdownRef}>
          <input
            id="batch-select"
            type="text"
            value={isBatchDropdownOpen ? batchSearchTerm : (batch || '')}
            onChange={(e) => {
              setBatchSearchTerm(e.target.value);
              setIsBatchDropdownOpen(true);
            }}
            onFocus={() => {
              setBatchSearchTerm('');
              setIsBatchDropdownOpen(true);
            }}
            disabled={loadingBatches}
            placeholder={loadingBatches ? 'Loading...' : 'Search batches...'}
            className="w-full px-6 py-5 text-xl bg-white dark:bg-[#3A3A3A] border border-[#2B2B2B]/15 dark:border-[#F5F1E8]/15 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#2B2B2B]/30 dark:focus:ring-[#F5F1E8]/30 transition-all"
          />

          {/* Batch Dropdown List */}
          {isBatchDropdownOpen && filteredBatches.length > 0 && (
            <div className="absolute z-10 w-full mt-2 bg-white dark:bg-[#2B2B2B] border border-[#2B2B2B]/15 dark:border-[#F5F1E8]/15 rounded-xl shadow-lg max-h-60 overflow-y-auto">
              {filteredBatches.map((batchOption) => (
                <button
                  key={batchOption}
                  onClick={() => selectBatch(batchOption)}
                  className="w-full px-6 py-4 text-xl text-left hover:bg-[#E8E4DC] dark:hover:bg-[#3A3A3A] transition-colors border-b border-[#2B2B2B]/10 dark:border-[#F5F1E8]/10 last:border-b-0"
                >
                  {batchOption}
                </button>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Add All Subjects Button */}
      {batch && (
        <button
          onClick={async () => {
            try {
              const response = await fetch(`/api/timetable/batches/${encodeURIComponent(batch)}/subjects`);
              if (!response.ok) throw new Error('Failed to fetch batch subjects');
              const data = await response.json();
              const newSubjects = (data.subjects as string[])
                .filter((code: string) => !subjects.some(s => s.code === code))
                .map((code: string) => ({ code }));
              if (newSubjects.length > 0) {
                setSubjects([...subjects, ...newSubjects]);
              }
            } catch (error) {
              console.error('Error fetching batch subjects:', error);
            }
          }}
          className="w-full px-6 py-4 text-lg font-medium bg-[#2B2B2B] dark:bg-[#F5F1E8] text-[#F5F1E8] dark:text-[#2B2B2B] rounded-xl hover:opacity-90 transition-all"
        >
          Add All Subjects from Batch
        </button>
      )}

      {/* Subjects Dropdown */}
      <div className="space-y-4">
        <label htmlFor="subject-search" className="block text-2xl">
          Subjects
        </label>
        <div className="relative" ref={dropdownRef}>
          <input
            id="subject-search"
            type="text"
            value={searchTerm}
            onChange={(e) => {
              setSearchTerm(e.target.value);
              setIsDropdownOpen(true);
            }}
            onFocus={() => setIsDropdownOpen(true)}
            disabled={loadingSubjects}
            placeholder={loadingSubjects ? 'Loading subjects...' : 'Search subject codes...'}
            className="w-full px-6 py-5 text-xl bg-white dark:bg-[#3A3A3A] border border-[#2B2B2B]/15 dark:border-[#F5F1E8]/15 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#2B2B2B]/30 dark:focus:ring-[#F5F1E8]/30 transition-all"
          />

          {/* Dropdown List */}
          {isDropdownOpen && filteredSubjects.length > 0 && (
            <div className="absolute z-10 w-full mt-2 bg-white dark:bg-[#2B2B2B] border border-[#2B2B2B]/15 dark:border-[#F5F1E8]/15 rounded-xl shadow-lg max-h-60 overflow-y-auto">
              {filteredSubjects.map((subjectCode) => (
                <button
                  key={subjectCode}
                  onClick={() => addSubject(subjectCode)}
                  className="w-full px-6 py-4 text-xl text-left hover:bg-[#E8E4DC] dark:hover:bg-[#3A3A3A] transition-colors border-b border-[#2B2B2B]/10 dark:border-[#F5F1E8]/10 last:border-b-0"
                >
                  {subjectCode}
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Selected Subjects List */}
        {subjects.length > 0 && (
          <div className="space-y-3 mt-6">
            {subjects.map((subject) => (
              <div
                key={subject.code}
                className="flex items-center justify-between gap-3 px-6 py-4 bg-white dark:bg-[#3A3A3A] border border-[#2B2B2B]/15 dark:border-[#F5F1E8]/15 rounded-xl group hover:border-[#2B2B2B]/30 dark:hover:border-[#F5F1E8]/30 transition-all"
              >
                <div className="flex-1 min-w-0">
                  {editingSubjectCode === subject.code ? (
                    <div className="flex items-center gap-2">
                      <span className="text-sm text-muted-foreground whitespace-nowrap">{subject.code}</span>
                      <span className="text-muted-foreground">→</span>
                      <input
                        type="text"
                        value={editingAlias}
                        onChange={(e) => {
                          if (e.target.value.length <= MAX_ALIAS_LENGTH) {
                            setEditingAlias(e.target.value);
                          }
                        }}
                        onKeyDown={(e) => {
                          if (e.key === 'Enter') {
                            saveAlias(subject.code);
                          } else if (e.key === 'Escape') {
                            setEditingSubjectCode(null);
                            setEditingAlias('');
                          }
                        }}
                        placeholder="Add nickname"
                        className="flex-1 px-3 py-1 text-lg bg-[#E8E4DC] dark:bg-[#2B2B2B] border border-[#2B2B2B]/15 dark:border-[#F5F1E8]/15 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#2B2B2B]/30 dark:focus:ring-[#F5F1E8]/30"
                        autoFocus
                      />
                      <span className="text-xs text-muted-foreground whitespace-nowrap">
                        {editingAlias.length}/{MAX_ALIAS_LENGTH}
                      </span>
                    </div>
                  ) : (
                    <div className="flex items-center gap-2">
                      {subject.alias ? (
                        <>
                          <span className="text-xl">{subject.alias}</span>
                          <span className="text-sm text-muted-foreground">({subject.code})</span>
                        </>
                      ) : (
                        <span className="text-xl">{subject.code}</span>
                      )}
                    </div>
                  )}
                </div>

                <div className="flex items-center gap-2">
                  {editingSubjectCode === subject.code ? (
                    <button
                      onClick={() => saveAlias(subject.code)}
                      className="p-2 hover:bg-[#E8E4DC] dark:hover:bg-[#2B2B2B] rounded-lg transition-all"
                      aria-label="Save alias"
                    >
                      <Check className="w-5 h-5" />
                    </button>
                  ) : (
                    <button
                      onClick={() => startEditingAlias(subject.code, subject.alias)}
                      className="p-2 hover:bg-[#E8E4DC] dark:hover:bg-[#2B2B2B] rounded-lg transition-all"
                      aria-label="Edit alias"
                    >
                      <Edit2 className="w-4 h-4" />
                    </button>
                  )}
                  <button
                    onClick={() => removeSubject(subject.code)}
                    className="p-2 hover:bg-[#E8E4DC] dark:hover:bg-[#2B2B2B] rounded-lg transition-all"
                    aria-label={`Remove ${subject.code}`}
                  >
                    <X className="w-5 h-5" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}