"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@/hooks/useUser";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface ServiceInstance {
  id: string;
  name: string;
  status: "running" | "stopped";
  health: "healthy" | "unhealthy";
  url: string;
  createdAt: string;
}

interface ServiceGroup {
  name: string;
  instances: ServiceInstance[];
  maxInstances: number;
}

export default function ServicesPage() {
  const [services, setServices] = useState<ServiceGroup[]>([]);
  const router = useRouter();
  const { isAdmin, isLoading } = useUser();

  useEffect(() => {
    if (!isAdmin && !isLoading) {
      router.push('/home');
    }
  }, [isAdmin, router, isLoading]);

  useEffect(() => {
    fetchServices();
  }, []);

  const fetchServices = async () => {
    try {
      const response = await fetch('/api/services');
      const data = await response.json();
      if (response.ok) {
        setServices(data.data);
      }
    } catch (error) {
      console.error('Error fetching services:', error);
    }
  };

  const addInstance = async (serviceName: string) => {
    try {
      const response = await fetch('/api/services', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ serviceName }),
      });

      if (response.ok) {
        fetchServices();
      }
    } catch (error) {
      console.error('Error adding instance:', error);
    }
  };

  const removeInstance = async (serviceName: string, instanceId: string) => {
    try {
      const response = await fetch(`/api/services/${instanceId}`, {
        method: 'DELETE',
      });

      if (response.ok) {
        fetchServices();
      }
    } catch (error) {
      console.error('Error removing instance:', error);
    }
  };

  if (!isAdmin) {
    return null;
  }

  return (
    <div className="container mx-auto py-8 px-4">
      <h1 className="text-3xl font-bold mb-8">Services Dashboard</h1>
      <div className="grid gap-6">
        {services.map((service) => (
          <Card key={service.name}>
            <CardHeader>
              <div className="flex justify-between items-center">
                <CardTitle>{service.name}</CardTitle>
                <Button
                  onClick={() => addInstance(service.name)}
                  disabled={true}
                >
                  Add Instance
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4">
                {service.instances.map((instance) => (
                  <div
                    key={instance.id}
                    className="flex items-center justify-between p-4 border rounded-lg"
                  >
                    <div className="flex items-center gap-4">
                      <div>
                        <h3 className="font-semibold">{instance.name}</h3>
                        <p className="text-sm text-gray-500">{instance.url}</p>
                      </div>
                      <Badge
                        variant={instance.health === "healthy" ? "default" : "destructive"}
                      >
                        {instance.health}
                      </Badge>
                      <Badge
                        variant={instance.status === "running" ? "default" : "secondary"}
                      >
                        {instance.status}
                      </Badge>
                    </div>
                    <Button
                      variant="destructive"
                      onClick={() => removeInstance(service.name, instance.id)}
                      disabled={true}
                    >
                      Remove
                    </Button>
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