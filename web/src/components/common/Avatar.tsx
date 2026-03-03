interface AvatarProps {
  name: string;
  url?: string;
  size?: 'sm' | 'md' | 'lg';
}

const colors = [
  'bg-blue-500', 'bg-green-500', 'bg-yellow-500', 'bg-red-500',
  'bg-purple-500', 'bg-pink-500', 'bg-indigo-500', 'bg-teal-500',
];

function getColor(name: string): string {
  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }
  return colors[Math.abs(hash) % colors.length];
}

export default function Avatar({ name, url, size = 'md' }: AvatarProps) {
  const dims = { sm: 'h-6 w-6 text-xs', md: 'h-9 w-9 text-sm', lg: 'h-12 w-12 text-lg' }[size];

  if (url) {
    return <img src={url} alt={name} className={`${dims} rounded-full object-cover flex-shrink-0`} />;
  }

  const initials = name
    .split(/[\s_-]+/)
    .map((w) => w[0])
    .slice(0, 2)
    .join('')
    .toUpperCase();

  return (
    <div className={`${dims} rounded-full flex-shrink-0 flex items-center justify-center text-white font-semibold ${getColor(name)}`}>
      {initials}
    </div>
  );
}
