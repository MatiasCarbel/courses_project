"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@/hooks/useUser";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface Port {
  ip: string;
  privatePort: number;
  publicPort: number;
  type: string;
}

interface Container {
  id: string;
  name: string;
  status: string;
  image: string;
  created: number;
  state: string;
  ports: Port[];
}

export default function ServicesPage() {
  const [containers, setContainers] = useState<Container[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const router = useRouter();
  const { isAdmin, isLoading } = useUser();

  useEffect(() => {
    if (!isLoading && !isAdmin) {
      router.push('/home');
      return;
    }

    const fetchContainers = async () => {
      try {
        const response = await fetch('/api/admin/containers', {
          credentials: 'include'
        });

        if (response.status === 401) {
          router.push('/login');
          return;
        }

        if (!response.ok) {
          throw new Error(`Error: ${response.statusText}`);
        }

        const data = await response.json();
        setContainers(data);
      } catch (err: any) {
        console.error('Error fetching containers:', err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    if (isAdmin) {
      fetchContainers();
      const interval = setInterval(fetchContainers, 5000);
      return () => clearInterval(interval);
    }
  }, [isAdmin, isLoading, router]);

  if (isLoading) return <div>Loading...</div>;
  if (!isAdmin) return null;
  if (error) return <div>Error: {error}</div>;

  const groupedContainers = containers.reduce((acc, container) => {
    const serviceName = container.name.split('-')[0];
    if (!acc[serviceName]) {
      acc[serviceName] = [];
    }
    acc[serviceName].push(container);
    return acc;
  }, {} as Record<string, Container[]>);

  return (
    <div className="container mx-auto py-8 px-4">
      <h1 className="text-3xl font-bold mb-8">Services Dashboard</h1>
      <div className="grid gap-6">
        {Object.entries(groupedContainers).map(([serviceName, serviceContainers]) => (
          <Card key={serviceName}>
            <CardHeader>
              <CardTitle className="capitalize">{serviceName} Service</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4">
                {serviceContainers.map((container) => (
                  <div
                    key={container.id}
                    className="flex items-center justify-between p-4 border rounded-lg"
                  >
                    <div className="flex items-center gap-4">
                      <div>
                        <h3 className="font-semibold">{container.name}</h3>
                        <p className="text-sm text-gray-500">ID: {container.id}</p>
                        <p className="text-sm text-gray-500">
                          Ports: {container.ports.map(p => `${p.publicPort}:${p.privatePort}`).join(', ')}
                        </p>
                      </div>
                      <Badge
                        variant={container.state === "running" ? "default" : "destructive"}
                      >
                        {container.state}
                      </Badge>
                    </div>
                    <div className="text-sm text-gray-500">
                      {new Date(container.created * 1000).toLocaleString()}
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
} 